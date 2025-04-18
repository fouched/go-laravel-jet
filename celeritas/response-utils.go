package celeritas

import (
	"github.com/fouched/toolkit/v2"
	"net/http"
)

func (c *Celeritas) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	var t toolkit.Tools
	return t.WriteJSON(w, status, data, headers...)
}
