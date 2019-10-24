package middleware

import (
	"fmt"
	"strings"
	"time"

	"iris-ticket/backend/config"
	"iris-ticket/backend/controllers"
	"iris-ticket/backend/middleware/jwts"
	"iris-ticket/backend/models"

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
	jj := JwtHandler().Get

	fmt.Printf("%v", jj)

	// auth
	if jj != nil {
		AuthToken(ctx)
	}
}

/**
 * 验证 jwt
 * @method JwtHandler
 */
func JwtHandler() *jwts.Middleware {
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
func AuthToken(ctx iris.Context) {
	u := ctx.Values().Get(config.Conf.Get("jwt.secert").(string)) //获取 token 信息
	if u == nil {
		ctx.StatusCode(iris.StatusBadRequest)
		_, _ = ctx.JSON(controllers.ApiResource(false, "Please take a token request", "error"))
		return
	} else {
		k := u.(*jwt.Token)
		token := models.GetOauthTokenByToken(k.Raw) //获取 access_token 信息
		if token.Revoked || token.ExpressIn < time.Now().Unix() {
			// ctx.StatusCode(http.StatusUnauthorized)
			ctx.JSON(controllers.ApiJson{Status: false, Data: "", Msg: "Token has expired"})
			return
		} else {
			ctx.Values().Set("auth_user_id", token.UserId)
		}
	}

	ctx.Next() // execute the "after" handler registered via `DoneGlobal`.
}

/**
return
	true:则跳过不需验证，如登录接口等...
	false:需要进一步验证
*/
func checkURL(reqPath string) bool {
	//config := iris.YAML("conf/app.yml")
	ignoreURLs := config.Conf.Get("server.ignore_urls").([]interface{})
	for _, v := range ignoreURLs {
		if reqPath == v {
			return true
		}
	}
	return false
}
