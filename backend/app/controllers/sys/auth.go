package sys

import (
	"iris-ticket/backend/app/config"
	"iris-ticket/backend/app/controllers/common"
	"iris-ticket/backend/app/models/db"
	"iris-ticket/backend/app/models/sys"
	
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jameskeane/bcrypt"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
)

type Auth struct{}

// 用户登录
func (Auth) Login(ctx iris.Context) {
	aul := sys.User{}
	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		common.ResErrSrv(ctx, err)
		return
	} else {
		if UserNameErr := common.Validate.Var(aul.Username, "required,min=3,max=20"); UserNameErr != nil {
			common.ResFail(ctx, "username format err")
			return
		} else if PwdErr := common.Validate.Var(aul.Password, "required,min=6,max=20"); PwdErr != nil {
			common.ResFail(ctx, "password format err")
		} else {
			ctx.StatusCode(iris.StatusOK)
			response, _, _ := CheckLogin(aul.Username, aul.Password)
			common.ResSuccess(ctx, response)
		}
	}
}

// 用户登出
func (Auth) Logout(ctx iris.Context) {
	// 删除uid
	uid, _ := ctx.Values().GetUint("auth_user_id")
	UserAdminLogout(uid)
	common.ResSuccess(ctx, uid)
}

/**
 * 校验用户登录
 * @method UserAdminCheckLogin
 * @param  {[type]}  username string [description]
 */
func UserAdminCheckLogin(username string) sys.User {
	aul := sys.User{}
	if err := db.DB.Where("username = ?", username).First(&aul).Error; err != nil {
		golog.Error("UserAdminCheckLoginErr ", err)
	}
	return aul
}

/**
* 用户退出登陆
* @method UserAdminLogout
* @param  {[type]} ids string [description]
 */
func UserAdminLogout(uid uint) string {
	ot := sys.UpdateOauthTokenByUserId(uid)
	return ot.Secret
}

/**
 * 判断用户是否登录
 * @method CheckLogin
 * @param  {[type]}  username string    [description]
 * @param  {[type]}  password string [description]
 */
func CheckLogin(username, password string) (response sys.Token, status bool, msg string) {
	user := UserAdminCheckLogin(username)
	if user.ID == 0 {
		msg = "user is not exist"
		return
	} else {
		if ok := bcrypt.Match(password, user.Password);ok {
			expireTime := time.Now().Add(time.Hour * time.Duration(config.Conf.Get("jwt.timeout").(int64))).Unix()
			jwtSecret := config.Conf.Get("jwt.secert").(string)
			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)
			claims["exp"] = expireTime
			claims["iat"] = time.Now().Unix()
			token.Claims = claims
			Tokenstring, err := token.SignedString([]byte(jwtSecret))

			if err != nil {
				msg = err.Error()
				return
			}

			oauthToken := new(sys.OauthToken)
			oauthToken.Token = Tokenstring
			oauthToken.UserId = user.ID
			oauthToken.Secret = jwtSecret
			oauthToken.Revoked = false
			oauthToken.ExpressIn = expireTime
			oauthToken.CreatedAt = time.Now()
			response = oauthToken.OauthTokenCreate()
			status = true
			msg = "success"

			return
		} else {
			msg = "password is error"
			return
		}
	}
}
