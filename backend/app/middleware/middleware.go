package middleware

import (
	"strings"
	"time"

	"github.com/kataras/golog"

	"iris-ticket/backend/app/config"
	"iris-ticket/backend/app/controllers/common"
	"iris-ticket/backend/app/middleware/casbins"
	"iris-ticket/backend/app/middleware/jwts"
	models "iris-ticket/backend/app/models/common"
	"iris-ticket/backend/app/models/sys"
	"iris-ticket/backend/app/pkg/convert"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

func ServeHTTP(ctx iris.Context) {
	path := ctx.Path()
	// 过滤静态资源、login接口、首页等...不需要验证
	if checkURL(path) || strings.Contains(path, "/statics") {
		ctx.Next()
		return
	}

	// jwt
	JwtMiddleware().Serve(ctx)

	// auth
	AuthTokenMiddleware(ctx)

	//casbin
	// skipperCasbinUrls := config.Conf.Get("server.skipper_casbin_urls").([]interface{})
	// CasbinMiddleware(AllowPathPrefixSkipper(skipperCasbinUrls...))
}

/**
return
	true:则跳过不需验证，如登录接口等...
	false:需要进一步验证
*/
func checkURL(reqPath string) bool {
	//config := iris.YAML("conf/app.yml")
	skipperTokenUrls := config.Conf.Get("server.skipper_token_urls").([]interface{})
	for _, v := range skipperTokenUrls {
		if reqPath == v {
			return true
		}
	}
	return false
}

/**
 * 验证 jwt
 * @method JwtHandler
 */
func JwtMiddleware() *jwts.Middleware {
	jwtSecert := config.Conf.Get("jwt.secert").(string)
	return jwts.New(jwts.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecert), nil
		},
		ContextKey:    jwtSecert,
		SigningMethod: jwt.SigningMethodHS256,
	})
}

/**
 * 判断 token 是否有效
 * 如果有效 就获取信息并且保存到请求里面
 * @method AuthToken
 * @param  {[type]}  ctx       iris.Context    [description]
 */
func AuthTokenMiddleware(ctx iris.Context) {
	u := ctx.Values().Get(config.Conf.Get("jwt.secert").(string)).(*jwt.Token) //获取 token 信息
	var model sys.OauthToken
	where := sys.OauthToken{}
	where.Token = u.Raw
	models.First(&where, &model)
	if model.Revoked || model.ExpressIn < time.Now().Unix() {
		common.ResFail(ctx, "Token has expired")
		return
	} else {
		ctx.Values().Set("auth_user_id", model.UserId)
	}
	// if everything is ok, next().
	ctx.Next()
}

// CasbinMiddleware casbin中间件
func CasbinMiddleware(skipper ...SkipperFunc) iris.Handler {
	return func(ctx iris.Context) {
		if len(skipper) > 0 && skipper[0](ctx) {
			ctx.Next()
			return
		}
		// 用户ID
		uid, err := ctx.Values().GetUint64("auth_user_id")
		if err != nil {
			common.ResFailCode(ctx, "token 无效3", 50008)
			return
		}
		if convert.ToUint64(uid) == common.SUPER_ADMIN_ID {
			ctx.Next()
			return
		}
		p := ctx.Path()
		m := ctx.Method()
		if b, err := casbins.CsbinCheckPermission(convert.ToString(uid), p, m); err != nil {
			common.ResFail(ctx, "err303"+err.Error())
			golog.Error("err303**")
			return
		} else if !b {
			common.ResFail(ctx, "没有访问权限")
			return
		}
		ctx.Next()
	}
}
