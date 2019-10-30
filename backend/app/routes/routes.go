package routes

import (
	sys "iris-ticket/backend/app/controllers/sys"
	"iris-ticket/backend/app/middleware"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/router"
)

func Register(app *iris.Application) {
	// allow cors
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	main := app.Party("/", crs).AllowMethods(iris.MethodOptions)

	main.Get("/", func(ctx iris.Context) { // 首页模块
		//_ = ctx.View("index.html")
		ctx.HTML("<h1 style='height: 1000px;line-height: 1000px;text-align: center;'>召唤师，欢迎来到王者峡谷</h1>")
	})

	app.RegisterView(iris.HTML("apidoc", ".html"))
	app.HandleDir("/apidoc", "../apidoc") // 设置静态资源

	api := app.Party("/api", crs).AllowMethods(iris.MethodOptions)
	api.Use(middleware.ServeHTTP)

	{
		auths := sys.Auth{}
		api.PartyFunc("/auth", func(auth router.Party) {
			auth.Post("/login", auths.Login)
			auth.Post("/logout", auths.Logout)
			auth.Post("/changepwd", auths.ChangePwd)
		})
		menus := sys.Menu{}
		api.PartyFunc("/menu", func(menu router.Party) {
			app.Get("/list", menus.List)
			app.Get("/detail", menus.Detail)
			app.Get("/allmenu", menus.AllMenu)
			app.Get("/menubuttonlist", menus.MenuButtonList)
			app.Post("/delete", menus.Delete)
			app.Post("/update", menus.Update)
			app.Post("/create", menus.Create)
		})
		users := sys.User{}
		api.PartyFunc("/user", func(user router.Party) {
			user.Get("/detail", users.Detail)
			user.Get("/list", users.List)
			user.Get("/detail", users.Detail)
			user.Get("/adminsroleidlist", users.AdminsRoleIDList)
			user.Post("/delete", users.Delete)
			user.Post("/update", users.Update)
			user.Post("/create", users.Create)
			user.Post("/setrole", users.SetRole)
		})
		roles := sys.Role{}
		api.PartyFunc("/role", func(role router.Party) {
			app.Get("/list", roles.List)
			app.Get("/detail", roles.Detail)
			app.Get("/rolemenuidlist", roles.RoleMenuIDList)
			app.Get("/allrole", roles.AllRole)
			app.Post("/delete", roles.Delete)
			app.Post("/update", roles.Update)
			app.Post("/create", roles.Create)
			app.Post("/setrole", roles.SetRole)
		})

	}
}
