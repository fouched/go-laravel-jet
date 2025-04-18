package main

import (
	"github.com/fouched/celeritas"
	"log"
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
	"os"
)

func initApplication() *application {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init celeritas
	cel := &celeritas.Celeritas{}
	err = cel.New(path)
	if err != nil {
		log.Fatal(err)
	}

	cel.AppName = "myapp"

	myMiddleware := &middleware.Middleware{
		App: cel,
	}

	myHandlers := &handlers.Handlers{App: cel}

	app := &application{
		App:        cel,
		Handlers:   myHandlers,
		Middleware: myMiddleware,
	}

	// set application routes to celeritas routes
	app.App.Routes = app.routes()

	// set the models
	app.Models = data.New(app.App.DB.Pool)
	myHandlers.Models = app.Models
	app.Middleware.Models = app.Models

	return app
}
