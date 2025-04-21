package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/fouched/celeritas"
	"myapp/data"
	"net/http"
)

func (h *Handlers) Form(w http.ResponseWriter, r *http.Request) {
	vars := make(jet.VarMap)
	validator := h.App.Validator(nil)
	vars.Set("validator", validator)
	vars.Set("user", data.User{})

	err := h.render(w, r, "form", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostForm(w http.ResponseWriter, r *http.Request) {
	// all form posts must be parsed
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	validator := h.App.Validator(nil)

	validator.Required(
		celeritas.Field{
			Name:  "first_name",
			Label: "First Name",
			Value: r.Form.Get("first_name"),
		},
		celeritas.Field{
			Name:  "last_name",
			Label: "Last Name",
			Value: r.Form.Get("last_name"),
		},
		celeritas.Field{
			Name:  "email",
			Label: "Email",
			Value: r.Form.Get("email"),
		},
	)

	validator.IsEmail(celeritas.Field{
		Name:  "email",
		Label: "Email",
		Value: r.Form.Get("email"),
	})

	validator.Check(len(r.Form.Get("first_name")) > 1, "first_name", "First Name must be at least two characters")
	//validator.Check(len(r.Form.Get("last_name")) > 1, "last_name", "Last Name must be at least two characters")
	// can also use IsLength for above
	validator.IsLength(
		celeritas.Field{
			Name:  "last_name",
			Label: "Last Name",
			Value: r.Form.Get("last_name"),
		}, 2)

	if validator.Valid() {
		fmt.Fprint(w, "valid data")
	} else {
		vars := make(jet.VarMap)
		vars.Set("validator", validator)

		var user data.User
		user.FirstName = r.Form.Get("first_name")
		user.LastName = r.Form.Get("last_name")
		user.Email = r.Form.Get("email")
		vars.Set("user", user)

		if err := h.render(w, r, "form", vars, nil); err != nil {
			h.App.ErrorLog.Println(err)
		}
	}

}
