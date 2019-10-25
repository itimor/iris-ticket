package controllers

import (
	"net/http"

	"iris-ticket/backend/models"

	"github.com/kataras/iris"
)

/**
 * @api {post} api/login UserLogin
 * @apiName UserLogin
 * @apiGroup auth
 *
 * @apiParamExample {json} Request-Example:
 *     {
 *       "username": "admin",
 *       "password": "admin"
 *     }
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": true,
 *       "msg": "sucess",
 *       "data": {
 *          "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzE4MTc5MTAsImlhdCI6MTU3MTgxNDMxMH0.5dAz2Fcfd1diaXzYONaehLB5tbf7Nyfa1HUGO3P4qew"
 *       }
 *    }
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 500 internal server error
 *     {
 *       "status": false,
 *       "msg": "parames err",
 *       "data": null
 *    }
 */
func UserLogin(ctx iris.Context) {
	aul := new(models.UserJson)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		_, _ = ctx.JSON(ApiResource(false, nil, "parames err"))
	} else {
		if UserNameErr := validate.Var(aul.Username, "required,min=4,max=20"); UserNameErr != nil {
			ctx.StatusCode(iris.StatusOK)
			_, _ = ctx.JSON(ApiResource(false, nil, "usrename format err"))
		} else if PwdErr := validate.Var(aul.Password, "required,min=5,max=20"); PwdErr != nil {
			ctx.StatusCode(iris.StatusOK)
			_, _ = ctx.JSON(ApiResource(false, nil, "password format err"))
		} else {
			ctx.StatusCode(iris.StatusOK)
			response, status, msg := models.CheckLogin(aul.Username, aul.Password)
			_, _ = ctx.JSON(ApiResource(status, response, msg))
		}
	}
}

/**
 * @api {get} api/logout UserLogout
 * @apiName UserLogout
 * @apiGroup auth
 *
 * @apiHeaderExample {json} Header-Example:
 *     {
 *       "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzE4MjkyMDgsImlhdCI6MTU3MTgyNTYwOH0.95H2py7V_zNiqjTxq4lV0Plfx1P32n6D-Lhmi1CoLOk"
 *     }
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": true,
 *       "msg": "sucess",
 *       "data": nul
 *    }
 *
 */
func UserLogout(ctx iris.Context) {
	// u := ctx.Values().Get("jwt").(*jwt.Token)   //获取 token 信息
	// token := models.GetOauthTokenByToken(u.Raw) //获取 access_token 信息
	// models.DeleteRequestTokenByToken(token)

	aui, _ := ctx.Values().GetUint("auth_user_id")
	models.UserAdminLogout(aui)

	ctx.StatusCode(http.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "success"))
}

/**
 * @api {post} api/changePasswd/:id UpdateUserPassword
 * @apiName UpdateUserPassword
 * @apiGroup auth
 *
 * @apiHeaderExample {json} Header-Example:
 *     {
 *       "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzE4MjkyMDgsImlhdCI6MTU3MTgyNTYwOH0.95H2py7V_zNiqjTxq4lV0Plfx1P32n6D-Lhmi1CoLOk"
 *     }
 *
 * @apiParamExample {json} Request-Example:
 *     {
 *       "password": "admin"
 *     }
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": true,
 *       "msg": "sucess",
 *       "data": : "$2a$10$z97tHB.IzV1zL3ENMPNcIua1GgHgNqkNHskiezZDYr525C"
 *    }
 *
 */
func UpdateUserPassword(ctx iris.Context) {
	aul := new(models.UserPassword)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(aul)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
		} else {
			id, _ := ctx.Params().GetInt("id")
			uid := uint(id)

			u, _ := models.UpdateUserPassword(aul.Password, uid)
			ctx.StatusCode(iris.StatusOK)
			if u.ID == 0 {
				_, _ = ctx.JSON(ApiResource(false, u.Password, "error"))
			} else {
				_, _ = ctx.JSON(ApiResource(true, u.Password, "success"))
			}
		}
	}
}
