package controllers

import (
	"iris-ticket/backend/app/controllers/common"
	models "iris-ticket/backend/app/models/common"
	"iris-ticket/backend/app/models/sys"

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
		if UserNameErr := common.Validate.Var(aul.UserName, "required,min=4,max=20"); UserNameErr != nil {
			common.ResFail(ctx, "username format err")
			return
		} else if PwdErr := common.Validate.Var(aul.Password, "required,min=5,max=20"); PwdErr != nil {
			common.ResFail(ctx, "password format err")
		} else {
			ctx.StatusCode(iris.StatusOK)
			response, status, msg := models.CheckLogin(aul.Username, aul.Password)
			_, _ = ctx.JSON(ApiResource(status, response, msg))
		}
	}
}
