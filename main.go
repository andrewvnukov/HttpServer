package main

import (
	"fmt"
	"net/http"
	"restapi/handler"
)

func main() {
	bookHandler := handler.NewBookHandler()
	userHandler := handler.NewUserHandler()
	storyHandler := handler.NewPurchaseHandler()

	http.Handle("/books/", http.StripPrefix("/books/", bookHandler))
	http.Handle("/users/", http.StripPrefix("/users/", userHandler))
	http.Handle("/story/", http.StripPrefix("/story/", storyHandler))

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
