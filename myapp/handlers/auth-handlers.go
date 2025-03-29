package handlers

import "net/http"

func (h *Handlers) UserLoginGet(w http.ResponseWriter, r *http.Request) {

	err := h.App.Render.Page(w, r, "login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) UserLoginPost(w http.ResponseWriter, r *http.Request) {

}
