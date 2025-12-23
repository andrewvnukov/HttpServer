package handler

import (
	"net/http"
	"restapi/model"
	"restapi/utils"
	"strconv"

	"github.com/gorilla/mux"
)

// UserHandler обработчик HTTP запросов для пользователей
// @Description Обработчик для работы с данными пользователей
type UserHandler struct {
	User model.UserHandler
}

// NewUserHandler создает новый экземпляр UserHandler
// @Summary Создать обработчик пользователей
// @Description Инициализирует и возвращает новый обработчик для работы с пользователями
// @Return http.Handler готовый обработчик HTTP запросов
func NewUserHandler() http.Handler {
	var u UserHandler
	u.User = model.UsersInit()
	return &u
}

// ServeHTTP обрабатывает входящие HTTP запросы для пользователей
// @Summary Основной обработчик запросов пользователей
// @Description Маршрутизирует запросы к соответствующим методам обработки пользователей
// @Router /users [get]
// @Router /users/{id} [get]
// @Router /users/add [post]
// @Router /users/update [post]
// @Router /users/{id} [delete]
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if idStr := mux.Vars(r)["id"]; idStr != "" {
			h.GetUser(w, r, idStr)
		} else {
			h.GetAllUsers(w, r)
		}
	case http.MethodPost:
		act := mux.Vars(r)["action"]
		switch act {
		case "add":
			h.AddUser(w, r)
		case "update":
			h.UpdateUser(w, r)
		}
	case http.MethodDelete:
		id, ok := mux.Vars(r)["id"]
		if ok {
			h.RemoveUser(w, r, id)
		} else {
			utils.ErrNotFoundApi(w, r)
		}
	}
}

// GetUser возвращает информацию о конкретном пользователе
// @Summary Получить пользователя по ID
// @Description Возвращает детальную информацию о пользователе по его идентификатору
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя" minimum(1)
// @Success 200 {object} model.User "Информация о пользователе"
// @Failure 404 {object} string "Пользователь не найден"
// @Failure 400 {object} string "Неверный ID"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, idStr string) {
	if id, err := strconv.Atoi(idStr); err != nil {
		utils.ErrNotFoundApi(w, r)
	} else {
		w.Header().Add("Content-Type", "JSON")
		w.Write(h.User.GetUser(id))
	}
}

// GetAllUsers возвращает список всех пользователей
// @Summary Получить всех пользователей
// @Description Возвращает полный список всех пользователей в системе
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} model.User "Список всех пользователей"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "JSON")
	w.Write(h.User.GetAllUsers())
}

// UpdateUser обновляет информацию о пользователе
// @Summary Обновить пользователя
// @Description Обновляет данные существующего пользователя
// @Tags users
// @Accept application/x-www-form-urlencoded
// @Produce plain
// @Param id formData int true "ID пользователя для обновления" minimum(1)
// @Param name formData string false "Новое имя пользователя" example("Иван")
// @Param surname formData string false "Новая фамилия пользователя" example("Иванов")
// @Success 200 {string} string "User updated successfully!"
// @Failure 404 {object} string "Пользователь не найден"
// @Failure 400 {object} string "Неверные данные запроса"
// @Failure 500 {object} string "Ошибка обновления"
// @Router /users/update [post]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := model.User{
		Name:    r.FormValue("name"),
		Surname: r.FormValue("surname"),
	}
	if id, err := strconv.Atoi(r.FormValue("id")); err != nil {
		utils.ErrNotFoundApi(w, r)
	} else {
		user.Id = id
		if err := h.User.UpdateUser(user); err != nil {
			utils.ErrUpdatingStorage(w, r)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("User updated successfully!"))
		}
	}
}

// AddUser добавляет нового пользователя
// @Summary Добавить нового пользователя
// @Description Создает нового пользователя в системе
// @Tags users
// @Accept application/x-www-form-urlencoded
// @Produce plain
// @Param name formData string true "Имя пользователя" example("Алексей")
// @Param surname formData string true "Фамилия пользователя" example("Петров")
// @Success 200 {string} string "User added successfully!"
// @Failure 400 {object} string "Неверные данные запроса"
// @Failure 409 {object} string "Пользователь уже существует"
// @Failure 500 {object} string "Ошибка добавления"
// @Router /users/add [post]
func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	newby := model.User{
		Name:    r.FormValue("name"),
		Surname: r.FormValue("surname"),
	}
	err := h.User.AddUser(newby)
	if err != nil {
		utils.ErrUpdatingStorage(w, r)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User added successfully!"))
	}
}

// RemoveUser удаляет пользователя
// @Summary Удалить пользователя
// @Description Удаляет пользователя из системы по его идентификатору
// @Tags users
// @Accept json
// @Produce plain
// @Param id path int true "ID пользователя для удаления" minimum(1)
// @Success 200 {string} string "User removed successfully!"
// @Failure 404 {object} string "Пользователь не найден"
// @Failure 500 {object} string "Ошибка удаления"
// @Router /users/{id} [delete]
func (h *UserHandler) RemoveUser(w http.ResponseWriter, r *http.Request, idStr string) {
	if id, err := strconv.Atoi(idStr); err != nil {
		utils.ErrNotFoundApi(w, r)
	} else {
		err := h.User.RemoveUser(id)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("User removed successfully!"))
		}
	}
}
