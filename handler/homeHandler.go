package handler

import (
	"net/http"
)

type HomeHandler struct {
	userHandler  http.Handler
	bookHandler  http.Handler
	storyHandler http.Handler
}

func NewHomeHandler() http.Handler {
	h := &HomeHandler{
		userHandler:  NewUserHandler(),
		bookHandler:  NewBookHandler(),
		storyHandler: NewPurchaseHandler(),
	}
	return h
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Handle("/books/", h.bookHandler)
}
