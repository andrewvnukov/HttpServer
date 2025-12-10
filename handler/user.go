package handler

import (
	"net/http"
	"restapi/model"
	"strconv"
)

type UserHandler struct {
	User model.UserHandler
}

func NewUserHandler() http.Handler {
	var u UserHandler
	u.User = model.UsersInit()
	return &u
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		switch r.URL.Path {
		case "getuser":
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.Header().Add("Content-Type", "application")
				w.Write(h.User.GetUser(id))
			}
		case "getcount":
			w.Write(h.User.GetCount())
		default:
			w.Write(h.User.GetAllUsers())
		}
	case "POST":
		if r.URL.Path == "adduser" {

		} else {
			w.Write([]byte("Not Found"))
		}
	case "DELETE":
		if r.URL.Path == "deleteuser" {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				err = h.User.RemoveUser(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}
		} else {
			w.Write([]byte("Not Found"))
		}
	}
}
