package handler

import (
	"net/http"
	"restapi/model"
	"strconv"
)

type BookHandler struct {
	Books model.Books
}

func NewBookHandler() http.Handler {
	h := &BookHandler{
		model.BooksInit(),
	}
	return h
}

func (h *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "get" {
			h.GetBook(w, r)
		} else {
			h.GetAllBooks(w, r)
		}
	case http.MethodPost:
		if r.URL.Path == "add" {
			h.AddBook(w, r)
		} else {
			h.UpdateBook(w, r)
		}

	case http.MethodDelete:
		if r.URL.Path == "remove" {
			h.RemoveBook(w, r)
		}
	}
}

func (h *BookHandler) AddBook(w http.ResponseWriter, r *http.Request) {
	book := model.BookModel{
		Name:   r.FormValue("name"),
		Author: r.FormValue("author"),
	}
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		price = 100
	}
	book.Price = price

	h.Books.AddBook(book)
	w.Write([]byte("Book added successfully"))
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	res, _ := strconv.Atoi(r.FormValue("id"))
	w.Write(h.Books.GetBook(res))
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Write(h.Books.GetAllBooks())
}

func (h *BookHandler) RemoveBook(w http.ResponseWriter, r *http.Request) {
	res, _ := strconv.Atoi(r.FormValue("id"))
	h.Books.RemoveBook(res)
	w.Write([]byte("Book removed successfully"))
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	book := model.BookModel{
		Name:   r.FormValue("name"),
		Author: r.FormValue("author"),
	}
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	if err != nil {
		price = 100
	}
	book.Price = price

	h.Books.UpdateBook(book)
	w.Write([]byte("Book updated successfully"))
}
