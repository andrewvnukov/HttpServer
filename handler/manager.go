package handler

import (
	"net/http"
)

type HandlerManager map[string]http.Handler

func NewHandlerManager() HandlerManager {
	return HandlerManager{
		"books": NewBookHandler(),
		"users": NewUserHandler(),
		"story": NewPurchaseHandler(),
	}
}
