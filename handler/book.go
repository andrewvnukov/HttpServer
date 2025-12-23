package handler

import (
	"net/http"
	"restapi/model"
	"restapi/utils"
	"strconv"

	_ "restapi/docs" // Импорт сгенерированной документации

	"github.com/gorilla/mux"
)

// BookHandler обработчик HTTP запросов для книг
// @Description Обработчик для работы с коллекцией книг
type BookHandler struct {
	Books model.Books
}

// NewBookHandler создает новый экземпляр BookHandler
// @Summary Создать обработчик книг
// @Description Инициализирует и возвращает новый обработчик для работы с книгами
// @Return http.Handler готовый обработчик HTTP запросов
func NewBookHandler() http.Handler {
	h := &BookHandler{
		model.BooksInit(),
	}
	return h
}

// ServeHTTP обрабатывает входящие HTTP запросы
// @Summary Основной обработчик запросов
// @Description Маршрутизирует запросы к соответствующим методам обработки
// @Param method path string true "HTTP метод"
// @Router /books [get]
// @Router /books/{id} [get]
// @Router /books/add [post]
// @Router /books/update [post]
// @Router /books/{id} [delete]
func (h *BookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if mux.Vars(r)["id"] == "" {
			h.GetAllBooks(w, r)
		} else {
			h.GetBook(w, r)
		}
	case http.MethodPost:
		if mux.Vars(r)["action"] == "add" {
			h.AddBook(w, r)
		} else if mux.Vars(r)["action"] == "update" {
			h.UpdateBook(w, r)
		} else {
			utils.ErrNotFoundApi(w, r)
		}

	case http.MethodDelete:
		if mux.Vars(r)["id"] == "" {
			utils.ErrNotFoundApi(w, r)
		} else {
			h.RemoveBook(w, r)
		}
	}
}

// AddBook добавляет новую книгу в коллекцию
// @Summary Добавить новую книгу
// @Description Создает новую книгу с указанными параметрами
// @Tags books
// @Accept application/x-www-form-urlencoded
// @Produce plain
// @Param name formData string true "Название книги" example("Война и мир")
// @Param author formData string true "Автор книги" example("Лев Толстой")
// @Param price formData number false "Цена книги" example(599.99)
// @Success 200 {string} string "Book added successfully"
// @Router /books/add [post]
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

// GetBook возвращает информацию о конкретной книге
// @Summary Получить книгу по ID
// @Description Возвращает детальную информацию о книге по её идентификатору
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "ID книги" minimum(1)
// @Success 200 {object} model.BookModel "Информация о книге"
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	res, _ := strconv.Atoi(mux.Vars(r)["id"])
	result := h.Books.GetBook(res)
	if result == nil {
		http.Error(w, "Книга не найдена", 400)
	} else {
		w.Write(result)
	}

}

// GetAllBooks возвращает список всех книг
// @Summary Получить все книги
// @Description Возвращает полный список книг в коллекции
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} model.BookModel "Список всех книг"
// @Router /books [get]
func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Write(h.Books.GetAllBooks())
}

// RemoveBook удаляет книгу из коллекции
// @Summary Удалить книгу
// @Description Удаляет книгу по указанному идентификатору
// @Tags books
// @Accept json
// @Produce plain
// @Param id path int true "ID книги для удаления" minimum(1)
// @Success 200 {string} string "Book removed successfully"
// @Router /books/{id} [delete]
func (h *BookHandler) RemoveBook(w http.ResponseWriter, r *http.Request) {
	res, _ := strconv.Atoi(mux.Vars(r)["id"])
	h.Books.RemoveBook(res)
	w.Write([]byte("Book removed successfully"))
}

// UpdateBook обновляет информацию о существующей книге
// @Summary Обновить книгу
// @Description Обновляет данные книги по её идентификатору
// @Tags books
// @Accept application/x-www-form-urlencoded
// @Produce plain
// @Param id formData int true "ID книги для обновления" minimum(1)
// @Param name formData string false "Новое название книги" example("Обновленное название")
// @Param author formData string false "Новый автор" example("Новый автор")
// @Param price formData number false "Новая цена" example(699.99)
// @Success 200 {string} string "Book updated successfully"
// @Router /books/update [post]
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

	fId, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Book id is required"))
	}
	book.Id = fId

	h.Books.UpdateBook(book)

	w.Write([]byte("Book updated successfully"))
}
