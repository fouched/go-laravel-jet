package middleware

import (
	"github.com/fouched/celeritas"
	"myapp/data"
)

type Middleware struct {
	App    *celeritas.Celeritas
	Models data.Models
}
