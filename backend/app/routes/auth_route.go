package routes

import (
	"iris-ticket/backend/app/controllers/sys"

	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
)

func AuthRoute(party iris.Party) {
	api := party.Party("/api")
	{
		auths := sys.Auth{}
		api.PartyFunc("/auth", func(auth router.Party) {
			auth.Post("/login", auths.Login)
			auth.Post("/logout", auths.Logout)
			auth.Post("/changepwd", auths.ChangePwd)
		})
	}
}
