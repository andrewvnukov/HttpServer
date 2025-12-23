package model

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"restapi/utils"
)

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type Users struct {
	Users []User `json:"users"`
	Total int    `json:"total"`
}

type UserHandler interface {
	Get() error
	Save() error
	AddUser(user User) error
	UpdateUser(user User) error
	RemoveUser(id int) error
	GetUser(id int) []byte
	GetAllUsers() []byte
	GetCount() []byte
}

func UsersInit() UserHandler {
	var u Users
	u.Get()
	return &u
}

func (u *Users) Get() error {
	data, err := os.ReadFile("./storage/users.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, &u); err != nil {
		if err == io.EOF {
			return nil
		} else {
			return err
		}
	}

	return nil
}

func (u *Users) Save() error {
	if data, err := json.Marshal(u); err != nil {
		return err
	} else {
		return os.WriteFile("./storage/users.json", data, 0644)
	}
}

func (u *Users) AddUser(user User) error {
	if u.Users == nil {
		user.Id = 0
		u.Users = []User{user}
		u.Total = 1
	} else {

		user.Id = u.Total
		u.Users = append(u.Users, user)
		u.Total++
	}

	err := u.Save()
	if err != nil {
		return nil
	} else {
		return err
	}
}

func (u *Users) UpdateUser(user User) error {
	for i, us := range u.Users {
		if us.Id == user.Id {
			if i == len(u.Users)-1 {
				u.Users = u.Users[:i]
			} else {
				u.Users = append(u.Users[:i], u.Users[i+1:]...)
			}
			return nil
		}
	}
	return errors.New("id not found")
}

func (u *Users) RemoveUser(id int) error {
	for i, user := range u.Users {
		if user.Id == id {
			if i == len(u.Users)-1 {
				u.Users = u.Users[:i]
			} else {
				u.Users = append(u.Users[:i], u.Users[i+1:]...)
			}
			u.Total--
			u.CheckUsers()
			return nil
		}
	}

	return errors.New("User not found")
}
func (u *Users) GetUser(id int) []byte {
	for _, user := range u.Users {
		if user.Id == id {
			return utils.MarshalThis(user)
		}
	}
	return nil
}
func (u *Users) GetAllUsers() []byte {
	return utils.MarshalThis(u)
}
func (u *Users) GetCount() []byte {
	return utils.MarshalThis(u.Total)
}
func (u *Users) CheckUsers() {
	for i := range u.Users {
		u.Users[i].Id = i
	}
}
