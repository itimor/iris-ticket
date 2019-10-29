package controllers

import (
	"iris-ticket/backend/app/controllers/common"
	models "iris-ticket/backend/app/models/common"
	"iris-ticket/backend/app/models/sys"

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

// 详情
func (User) Detail(ctx iris.Context) {
	id := common.GetQueryToUint64(ctx, "id")
	var model sys.User
	where := sys.User{}
	where.ID = id
	_, err := models.First(&where, &model)
	if err != nil {
		common.ResErrSrv(ctx, err)
		return
	}
	model.Password = ""
	common.ResSuccess(ctx, &model)
}
