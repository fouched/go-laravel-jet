package main

import (
	"github.com/fouched/celeritas"
	"log"
	"myapp/data"
	"myapp/handlers"
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

	myHandlers := &handlers.Handlers{App: cel}

	app := &application{
		App:      cel,
		Handlers: myHandlers,
	}

	// set application routes to celeritas routes
	app.App.Routes = app.routes()
	app.Models = data.New(app.App.DB.Pool)

	return app
}
