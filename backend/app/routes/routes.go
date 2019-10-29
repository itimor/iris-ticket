package routes

import (
	sys "iris-ticket/backend/app/controllers/sys"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
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

	api := app.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	//v1.Use(middleware.ServeHTTP)

	{
		users := sys.User{}
		api.Get("/{id:uint}", users.List)
	}
}
