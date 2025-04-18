package main

import (
	"github.com/fouched/celeritas"
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
)

type application struct {
	App        *celeritas.Celeritas
	Handlers   *handlers.Handlers
	Models     data.Models
	Middleware *middleware.Middleware
}

func main() {
	c := initApplication()
	c.App.ListenAndServe()
}
