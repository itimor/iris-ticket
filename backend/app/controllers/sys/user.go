package controllers

import (
	"github.com/kataras/iris"
)

type ApiJson struct {
	Status bool        `json:"status"`
	Msg    interface{} `json:"msg"`
	Data   interface{} `json:"data"`
}

func ApiResource(status bool, objects interface{}, msg string) (apijson *ApiJson) {
	apijson = &ApiJson{Status: status, Data: objects, Msg: msg}
	return
}

type User struct{}

func (User) List(ctx iris.Context) {
	id, _ := ctx.Params().GetUint64("id")
	// user, _ := sys.GetUserById(id)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, id, "success"))
}
