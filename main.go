// file: main.go

package main

import (
	"github.com/iris-contrib/middleware/cors"
	"github.com/jlb0906/mymovie/controller"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	// Load the template files.
	//app.RegisterView(iris.HTML("view", ".html"))
	// cors
	app.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
		MaxAge:           600,
		AllowedMethods:   []string{iris.MethodGet, iris.MethodPost, iris.MethodOptions, iris.MethodHead, iris.MethodDelete, iris.MethodPut},
		AllowedHeaders:   []string{"*"},
	}))
	app.AllowMethods(iris.MethodOptions)
	// Serve our controller.
	mvc.Configure(app.Party("/movie"), movies)
	// use swagger middleware to
	app.Run(
		// Start the web server at localhost:8080
		iris.Addr(":8081"),
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)
}

// note the mvc.Application, it's not iris.Application.
func movies(app *mvc.Application) {
	// Add the basic authentication(admin:password) middleware
	//app.Router.Use(middleware.BasicAuth)
	// serve our movies controller.
	app.Handle(new(controller.MovieController))
}
