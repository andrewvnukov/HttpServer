package utils

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	ErrNotFoundApi       = http.HandlerFunc(ErrNotFoundApiFunc)
	ErrNotFoundApiCustom = http.HandlerFunc(ErrNotFoundApiCustomFunc)
	ErrUpdatingStorage   = http.HandlerFunc(ErrUpdatingStorageApiFunc)
)

func ErrNotFoundApiFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s\nNot Found", mux.Vars(r)["version"])
}

func ErrNotFoundApiCustomFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Sory, but %s was'nt found", mux.Vars(r)["version"])
}

func ErrUpdatingStorageApiFunc(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Got error while updating storage")
}
