package routes

import (
	sys "iris-ticket/backend/app/controllers/sys"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func Register(api *iris.Application) {
	// allow cors
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
	})

	main := api.Party("/", crs).AllowMethods(iris.MethodOptions)

	main.Get("/", func(ctx iris.Context) { // 首页模块
		//_ = ctx.View("index.html")
		ctx.HTML("<h1 style='height: 1000px;line-height: 1000px;text-align: center;'>召唤师，欢迎来到王者峡谷</h1>")
	})

	v1 := api.Party("/v1", crs).AllowMethods(iris.MethodOptions)
	//v1.Use(middleware.ServeHTTP)
	// v1.Use(middleware.ServeHTTP, middleware.JwtHandler().Serve, middleware.AuthToken)

	{
		users := sys.User{}
		v1.Get("/{id:uint}", users.List)

		// v1.PartyFunc("/api", func(admin router.Party) {
		// 	admin.PartyFunc("/auth", func(auth router.Party) {
		// 		auth.Post("/login", controllers.UserLogin)
		// 		auth.Get("/logout", controllers.UserLogout)
		// 		auth.Patch("/changePasswd/{id:uint}", controllers.UpdateUserPassword)
		// 	})
		// 	admin.PartyFunc("/users", func(users router.Party) {
		// 		users.Get("/", controllers.GetAllUsers)
		// 		users.Get("/{id:uint}", controllers.GetUser)
		// 		users.Post("/", controllers.CreateUser)
		// 		users.Put("/{id:uint}", controllers.UpdateUser)
		// 		users.Delete("/{id:uint}", controllers.DeleteUser)
		// 		users.Get("/profile", controllers.GetProfile)
		// 	})
		// 	admin.PartyFunc("/roles", func(roles router.Party) {
		// 		roles.Get("/", controllers.GetAllRoles)
		// 		roles.Get("/{id:uint}", controllers.GetRole)
		// 		roles.Post("/", controllers.CreateRole)
		// 		roles.Put("/{id:uint}", controllers.UpdateRole)
		// 		roles.Delete("/{id:uint}", controllers.DeleteRole)
		// 	})
		// })
	}
}
