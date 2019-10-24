package controllers

import (
	"fmt"

	"iris-ticket/backend/models"
	"iris-ticket/backend/utils"

	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
)

/**
 * @api {get} api/roles/:id GetRole
 * @apiName GetRole
 * @apiGroup roles
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": true,
 *       "msg": "sucess",
 *       "data": {
 *          "Name": "test"
 *       }
 *    }
 *
 */
func GetRole(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	role, _ := models.GetRoleById(id)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, role, "success"))
}

/**
 * @api {post} api/roles CreateRole
 * @apiName CreateRole
 * @apiGroup roles
 *
 * @apiParamExample {json} Request-Example:
 *     {
 *       "name": "name",
 *       "display_name": "display_name",
 *       "desc": "desc"
 *     }
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": true,
 *       "msg": "sucess",
 *       "data": {
 *          "Name": "test",
 *          "DisplayName": "测试",
 *          "desc": "test",
 *       }
 *    }
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 500 internal server error
 *     {
 *       "status": false,
 *       "msg": "error",
 *       "data": null
 *    }
 */
func CreateRole(ctx iris.Context) {
	role := new(models.RoleJson)

	if err := ctx.ReadJSON(&role); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(role)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println()
				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.Type())
				fmt.Println(err.Param())
				fmt.Println()
			}
		} else {
			u, _ := models.CreateRole(role)
			if u.ID == 0 {
				ctx.StatusCode(iris.StatusInternalServerError)
				_, _ = ctx.JSON(ApiResource(false, nil, "error"))
			} else {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(true, u, "success"))
			}
		}
	}
}

func UpdateRole(ctx iris.Context) {

	role := new(models.RoleJson)

	if err := ctx.ReadJSON(&role); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(role)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			for _, err := range err.(validator.ValidationErrors) {
				fmt.Println()
				fmt.Println(err.Namespace())
				fmt.Println(err.Field())
				fmt.Println(err.Type())
				fmt.Println(err.Param())
				fmt.Println()
			}
		} else {
			id, _ := ctx.Params().GetInt("id")
			uid := uint(id)

			u, _ := models.UpdateRole(role, uid)
			ctx.StatusCode(iris.StatusOK)
			if u.ID == 0 {
				_, _ = ctx.JSON(ApiResource(false, u, "error"))
			} else {
				_, _ = ctx.JSON(ApiResource(true, u, "success"))
			}
		}
	}
}

func DeleteRole(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	models.DeleteRoleById(id)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "success"))
}

func GetAllRoles(ctx iris.Context) {
	offset := utils.Tool.ParseInt(ctx.FormValue("offset"), 1)
	limit := utils.Tool.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	orderBy := ctx.FormValue("orderBy")

	roles, _ := models.GetAllRoles(name, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, roles, "success"))
}
