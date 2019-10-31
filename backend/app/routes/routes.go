package routes

import (
	"iris-ticket/backend/app/controllers/sys"
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
			menu.Get("/list", menus.List)
			menu.Get("/detail", menus.Detail)
			menu.Get("/allmenu", menus.AllMenu)
			menu.Get("/menubuttonlist", menus.MenuButtonList)
			menu.Post("/delete", menus.Delete)
			menu.Post("/update", menus.Update)
			menu.Post("/create", menus.Create)
		})
		users := sys.User{}
		api.PartyFunc("/user", func(user router.Party) {
			user.Get("/detail", users.Detail)
			user.Get("/list", users.List)
			user.Get("/userroleidlist", users.AdminsRoleIDList)
			user.Post("/delete", users.Delete)
			user.Post("/update", users.Update)
			user.Post("/create", users.Create)
			user.Post("/setrole", users.SetRole)
		})
		roles := sys.Role{}
		api.PartyFunc("/role", func(role router.Party) {
			role.Get("/list", roles.List)
			role.Get("/detail", roles.Detail)
			role.Get("/rolemenuidlist", roles.RoleMenuIDList)
			role.Get("/allrole", roles.AllRole)
			role.Post("/delete", roles.Delete)
			role.Post("/update", roles.Update)
			role.Post("/create", roles.Create)
			role.Post("/setrole", roles.SetRole)
		})
	}
}
