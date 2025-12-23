package model

import (
	"encoding/json"
	"os"
	"restapi/utils"

	_ "restapi/docs" // Импорт сгенерированной документации
)

// BookModel представляет модель книги
// @Description Информация о книге
type BookModel struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

// Library представляет библиотеку книг
// @Description Информация о библиотеке
type Library struct {
	Books      []BookModel `json:"books"`
	TotalBooks int         `json:"total"`
}

// Books представляет интерфейс для работы с книгами
type Books interface {
	Get() error
	Save() error
	AddBook(book BookModel)
	RemoveBook(id int)
	UpdateBook(book BookModel)
	GetBook(id int) []byte
	GetAllBooks() []byte
	GetCount() int
}

func BooksInit() Books {
	var l Library
	err := l.Get()
	if err != nil {
		panic(err)
	}
	return &l
}

func (l *Library) Get() error {
	data, err := os.ReadFile("./storage/books.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	err = json.Unmarshal(data, &l)
	if err != nil {
		return err
	}

	return nil
}
func (l *Library) GetBook(id int) []byte {
	for _, book := range l.Books {
		if book.Id == id {
			return utils.MarshalThis(book)
		}
	}
	return nil
}
func (l *Library) GetAllBooks() []byte {
	return utils.MarshalThis(l)
}
func (l *Library) GetCount() int {
	return l.TotalBooks
}
func (l *Library) Save() error {
	data, err := json.Marshal(l)
	if err != nil {
		return err
	}

	if err := os.WriteFile("./storage/books.json", data, 0644); err != nil {
		return err
	}

	return nil
}
func (l *Library) AddBook(book BookModel) {
	book.Id = l.TotalBooks + 1
	l.Books = append(l.Books, book)
	l.TotalBooks++
	l.Save()
}
func (l *Library) RemoveBook(id int) {
	for i, book := range l.Books {
		if book.Id == id {
			if i == len(l.Books)-1 {
				l.Books = l.Books[:i]
			} else {
				l.Books = append(l.Books[:i], l.Books[i+1:]...)
			}
		}
	}
	l.TotalBooks--
	l.Save()
}
func (l *Library) UpdateBook(book BookModel) {
	for i, b := range l.Books {
		if b.Id == book.Id {
			l.Books[i] = book
			l.Save()
		}
	}
}
