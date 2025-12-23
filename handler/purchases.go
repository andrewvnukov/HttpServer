package handler

import (
	"fmt"
	"net/http"
	"restapi/model"
	"restapi/utils"
	"strconv"

	"github.com/gorilla/mux"
)

// PurchaseHandler обработчик HTTP запросов для истории покупок/аренды
// @Description Обработчик для работы с историей покупок или аренды книг
type PurchaseHandler struct {
	Purchase model.StoryHandler
}

// NewPurchaseHandler создает новый экземпляр PurchaseHandler
// @Summary Создать обработчик истории покупок
// @Description Инициализирует и возвращает новый обработчик для работы с историей покупок/аренды
// @Return http.Handler готовый обработчик HTTP запросов
func NewPurchaseHandler() http.Handler {
	var p PurchaseHandler
	p.Purchase = model.StoryInit()

	return &p
}

// ServeHTTP обрабатывает входящие HTTP запросы для истории покупок
// @Summary Основной обработчик запросов истории покупок
// @Description Маршрутизирует запросы к соответствующим методам обработки истории покупок
// @Router /story [get]
// @Router /story/id/{id} [get]
// @Router /story/book/{id} [get]
// @Router /story/user/{id} [get]
// @Router /story [post]
// @Router /story/update/{id} [put]
// @Router /story/endpurchase/{id} [put]
// @Router /story/id/{id} [delete]
// @Router /story/book/{id} [delete]
// @Router /story/user/{id} [delete]
func (h *PurchaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// GET запросы
	case http.MethodGet:
		idStr, ok := mux.Vars(r)["id"]
		if !ok {
			if idStr == "" {
				w.Write(h.Purchase.GetAll())
				return
			}
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		switch mux.Vars(r)["action"] {
		case "id":
			w.Header().Add("Content-Type", "Json")
			w.Write(h.Purchase.GetById(id))
		case "book":
			w.Header().Add("Content-Type", "Json")
			w.Write(h.Purchase.GetByBook(id))
		case "user":
			w.Header().Add("Content-Type", "Json")
			w.Write(h.Purchase.GetByUser(id))
		default:
			utils.ErrNotFoundApi(w, r)
		}
	// POST запросы
	case http.MethodPost:
		h.AddPurchase(w, r)
	// PUT запросы
	case http.MethodPut:
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utils.ErrNotFoundApi(w, r)
		}
		switch mux.Vars(r)["action"] {
		case "update":
			h.UpdatePurchase(w, r, id)
		case "endpurchase":
			err := h.Purchase.EndPurchase(id)
			if err != nil {
				utils.ErrUpdatingStorage(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Purchase ended successfully!"))
			}
		default:
			utils.ErrNotFoundApi(w, r)
		}
	// DELETE запросы
	case http.MethodDelete:
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			utils.ErrNotFoundApi(w, r)
		}
		switch mux.Vars(r)["action"] {
		case "id":
			err := h.Purchase.DelPurchase(id)
			if err != nil {
				utils.ErrUpdatingStorage(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Purchase deleted successfully!"))
			}
		case "book":
			err := h.Purchase.DelPurchaseByBook(id)
			if err != nil {
				utils.ErrUpdatingStorage(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Purchase deleted successfully!"))
			}
		case "user":
			err := h.Purchase.DelPurchaseByUser(id)
			if err != nil {
				utils.ErrUpdatingStorage(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Purchase deleted successfully!"))
			}
		}
	}
}

// AddPurchase добавляет новую запись о покупке/аренде
// @Summary Добавить новую покупку/аренду
// @Description Создает новую запись о покупке или аренде книги пользователем
// @Tags purchases
// @Accept application/x-www-form-urlencoded
// @Produce plain
// @Param book_id formData int true "ID книги" example(1)
// @Param user_id formData int true "ID пользователя" example(1)
// @Success 200 {string} string "Loan succesfully added"
// @Failure 400 {object} string "Неверные данные запроса"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /story [post]
func (h *PurchaseHandler) AddPurchase(w http.ResponseWriter, r *http.Request) {
	bookIdStr, userIdStr := r.FormValue("book_id"), r.FormValue("user_id")
	bookId, err := strconv.Atoi(bookIdStr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ошибка при обработке айди книги: %s", err.Error())
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ошибка при обработке айди пользователя: %s", err.Error())
		return
	}
	loan := model.Purchase{
		BookId: bookId,
		UserId: userId,
	}
	err = h.Purchase.AddPurchase(loan)
	if err != nil {
		utils.ErrUpdatingStorage(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Loan succesfully added"))
	}
}

// UpdatePurchase обновляет информацию о покупке/аренде
// @Summary Обновить информацию о покупке
// @Description Обновляет данные существующей записи о покупке или аренде
// @Tags purchases
// @Accept application/x-www-form-urlencoded
// @Produce plain
// @Param id path int true "ID записи о покупке" example(1)
// @Param book_id formData int false "Новый ID книги" example(2)
// @Param user_id formData int false "Новый ID пользователя" example(3)
// @Success 200 {string} string "Loan succesfully updated"
// @Failure 404 {object} string "Запись не найдена"
// @Failure 400 {object} string "Неверные данные запроса"
// @Failure 500 {object} string "Внутренняя ошибка сервера"
// @Router /story/update/{id} [put]
func (h *PurchaseHandler) UpdatePurchase(w http.ResponseWriter, r *http.Request, id int) {
	bookIdStr, userIdStr := r.FormValue("book_id"), r.FormValue("user_id")
	bookId, err := strconv.Atoi(bookIdStr)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ошибка при обработке айди книги: %s", err.Error())
		return
	}
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ошибка при обработке айди пользователя: %s", err.Error())
		return
	}
	loan := model.Purchase{
		Id:     id,
		BookId: bookId,
		UserId: userId,
	}
	err = h.Purchase.UpdatePurchase(loan)
	if err != nil {
		utils.ErrUpdatingStorage(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Loan succesfully updated"))
	}
}

// Также добавьте эти методы с аннотациями (если они есть в вашем коде):

// GetAllPurchases возвращает все записи о покупках
// @Summary Получить все покупки
// @Description Возвращает полный список всех записей о покупках/аренде
// @Tags purchases
// @Accept json
// @Produce json
// @Success 200 {array} model.Purchase "Список всех покупок"
// @Router /story [get]
// Примечание: Этот метод обрабатывается в ServeHTTP

// GetPurchaseById возвращает запись о покупке по ID
// @Summary Получить покупку по ID
// @Description Возвращает информацию о конкретной покупке по её идентификатору
// @Tags purchases
// @Accept json
// @Produce json
// @Param id path int true "ID покупки" example(1)
// @Success 200 {object} model.Purchase "Информация о покупке"
// @Failure 404 {object} string "Покупка не найдена"
// @Router /story/id/{id} [get]
// Примечание: Этот метод обрабатывается в ServeHTTP

// GetPurchasesByBook возвращает покупки по ID книги
// @Summary Получить покупки по книге
// @Description Возвращает все записи о покупках/аренде для указанной книги
// @Tags purchases
// @Accept json
// @Produce json
// @Param id path int true "ID книги" example(1)
// @Success 200 {array} model.Purchase "Список покупок для книги"
// @Failure 404 {object} string "Книга не найдена"
// @Router /story/book/{id} [get]
// Примечание: Этот метод обрабатывается в ServeHTTP

// GetPurchasesByUser возвращает покупки по ID пользователя
// @Summary Получить покупки по пользователю
// @Description Возвращает все записи о покупках/аренде для указанного пользователя
// @Tags purchases
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя" example(1)
// @Success 200 {array} model.Purchase "Список покупок пользователя"
// @Failure 404 {object} string "Пользователь не найдена"
// @Router /story/user/{id} [get]
// Примечание: Этот метод обрабатывается в ServeHTTP

// EndPurchase завершает покупку/аренду
// @Summary Завершить покупку
// @Description Отмечает покупку/аренду как завершенную
// @Tags purchases
// @Accept json
// @Produce plain
// @Param id path int true "ID покупки" example(1)
// @Success 200 {string} string "Purchase ended successfully!"
// @Failure 404 {object} string "Покупка не найдена"
// @Failure 500 {object} string "Ошибка обновления"
// @Router /story/endpurchase/{id} [put]
// Примечание: Этот метод обрабатывается в ServeHTTP

// DeletePurchase удаляет запись о покупке по ID
// @Summary Удалить покупку по ID
// @Description Удаляет запись о покупке по её идентификатору
// @Tags purchases
// @Accept json
// @Produce plain
// @Param id path int true "ID покупки" example(1)
// @Success 200 {string} string "Purchase deleted successfully!"
// @Failure 404 {object} string "Покупка не найдена"
// @Failure 500 {object} string "Ошибка удаления"
// @Router /story/id/{id} [delete]
// Примечание: Этот метод обрабатывается в ServeHTTP

// DeletePurchaseByBook удаляет покупки по ID книги
// @Summary Удалить покупки по книге
// @Description Удаляет все записи о покупках/аренде для указанной книги
// @Tags purchases
// @Accept json
// @Produce plain
// @Param id path int true "ID книги" example(1)
// @Success 200 {string} string "Purchase deleted successfully!"
// @Failure 404 {object} string "Книга не найдена"
// @Failure 500 {object} string "Ошибка удаления"
// @Router /story/book/{id} [delete]
// Примечание: Этот метод обрабатывается в ServeHTTP

// DeletePurchaseByUser удаляет покупки по ID пользователя
// @Summary Удалить покупки по пользователю
// @Description Удаляет все записи о покупках/аренде для указанного пользователя
// @Tags purchases
// @Accept json
// @Produce plain
// @Param id path int true "ID пользователя" example(1)
// @Success 200 {string} string "Purchase deleted successfully!"
// @Failure 404 {object} string "Пользователь не найдена"
// @Failure 500 {object} string "Ошибка удаления"
// @Router /story/user/{id} [delete]
// Примечание: Этот метод обрабатывается в ServeHTTP
