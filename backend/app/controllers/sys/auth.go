package sys

import (
	"time"

	"iris-ticket/backend/app/config"
	"iris-ticket/backend/app/controllers/common"
	models "iris-ticket/backend/app/models/common"
	"iris-ticket/backend/app/models/sys"

	"github.com/dgrijalva/jwt-go"
	"github.com/jameskeane/bcrypt"
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
			if response, status, _ := CheckLogin(ctx, aul.Username, aul.Password); status {
				common.ResSuccess(ctx, response)
			}
		}
	}
}

// 用户登出
func (Auth) Logout(ctx iris.Context) {
	// 删除uid
	uid, _ := ctx.Values().GetUint64("auth_user_id")
	where := sys.OauthToken{}
	where.Revoked = false
	where.UserId = uid
	modelOld := sys.OauthToken{}
	_, err := models.First(&where, &modelOld)
	if err != nil {
		common.ResErrSrv(ctx, err)
		return
	}
	modelNew := sys.OauthToken{Revoked: true}
	err = models.Updates(&modelOld, &modelNew)
	if err != nil {
		common.ResErrSrv(ctx, err)
		return
	}
	common.ResSuccess(ctx, uid)
}

// 用户修改密码
func (Auth) ChangePwd(ctx iris.Context) {
	modelNew := sys.User{}
	if err := ctx.ReadJSON(&modelNew); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		common.ResErrSrv(ctx, err)
		return
	} else {
		uid, _ := ctx.Values().GetUint64("auth_user_id")
		where := sys.User{}
		where.ID = uid
		modelOld := sys.User{}
		_, err := models.First(&where, &modelOld)
		if err != nil {
			common.ResErrSrv(ctx, err)
			return
		}
		salt, _ := bcrypt.Salt(10)
		hash, _ := bcrypt.Hash(modelNew.Password, salt)
		modelNew.Password = string(hash)
		err = models.Updates(&modelOld, &modelNew)
		if err != nil {
			common.ResErrSrv(ctx, err)
			return
		}
		common.ResSuccess(ctx, "password change success")
	}
}

/**
 * 校验用户登录
 * @method UserAdminCheckLogin
 * @param  {[type]}  username string [description]
 */
func UserAdminCheckLogin(ctx iris.Context, username string) (model sys.User) {
	where := sys.User{}
	where.Username = username
	model = sys.User{}
	_, err := models.First(&where, &model)
	if err != nil {
		common.ResFail(ctx, "操作失败")
		return
	}
	return
}

/**
 * 判断用户是否登录
 * @method CheckLogin
 * @param  {[type]}  username string    [description]
 * @param  {[type]}  password string [description]
 */
func CheckLogin(ctx iris.Context, username, password string) (response string, status bool, msg string) {
	user := UserAdminCheckLogin(ctx, username)
	if user.ID == 0 {
		status = false
		msg = "user is not exist"
		return
	} else {
		if ok := bcrypt.Match(password, user.Password); ok {
			expireTime := time.Now().Add(time.Hour * time.Duration(config.Conf.Get("jwt.timeout").(int64))).Unix()
			jwtSecret := config.Conf.Get("jwt.secert").(string)
			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)
			claims["exp"] = expireTime
			claims["iat"] = time.Now().Unix()
			token.Claims = claims
			Tokenstring, err := token.SignedString([]byte(jwtSecret))

			if err != nil {
				common.ResFail(ctx, err.Error())
				return
			}

			oauthToken := new(sys.OauthToken)
			oauthToken.Token = Tokenstring
			oauthToken.UserId = user.ID
			oauthToken.Secret = jwtSecret
			oauthToken.Revoked = false
			oauthToken.ExpressIn = expireTime
			oauthToken.CreatedAt = time.Now()
			err = models.Create(&oauthToken)
			if err != nil {
				status = false
			} else {
				response = Tokenstring
				status = true
			}
			return
		} else {
			common.ResFail(ctx, "密码错误")
			return
		}
	}
}
