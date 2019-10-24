package middleware

import (
	"time"

	"iris-ticket/backend/controllers"
	"iris-ticket/backend/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

/**
 * 判断 token 是否有效
 * 如果有效 就获取信息并且保存到请求里面
 * @method AuthToken
 * @param  {[type]}  ctx       iris.Context    [description]
 */
func AuthToken(ctx iris.Context) {
	u := ctx.Values().Get("jwt").(*jwt.Token)   //获取 token 信息
	token := models.GetOauthTokenByToken(u.Raw) //获取 access_token 信息
	if token.Revoked || token.ExpressIn < time.Now().Unix() {
		// ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(controllers.ApiJson{Status: false, Data: "", Msg: "Token has expired"})
		return
	} else {
		ctx.Values().Set("auth_user_id", token.UserId)
	}

	ctx.Next() // execute the "after" handler registered via `DoneGlobal`.
}
