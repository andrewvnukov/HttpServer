package model

import (
	"encoding/json"
	"os"
	"restapi/utils"
	"time"
)

type Purchase struct {
	Id     int       `json:"id"`
	BookId int       `json:"book_id"`
	UserId int       `json:"user_id"`
	TookAt time.Time `json:"start_at"`
	EndAt  time.Time `json:"end_at"`
}

type Story struct {
	Purchases []Purchase `json:"purchases"`
	LastId    int        `json:"lastid"`
}

type StoryHandler interface {
	Get()
	Save()
	GetAll() []byte
	GetByUser(int) []byte
	GetByBook(int) []byte
	GetById(int) []byte
	AddPurchase(Purchase)
	EndPurchase(int)
	DelPurchase(int)
	UpdatePurchase(Purchase)
}

func StoryInit() StoryHandler {
	var s Story
	s.Get()
	return &s
}
func (s *Story) Get() {
	if file, err := os.ReadFile("./storage/purchases.json"); err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	} else {
		err = json.Unmarshal(file, &s)
		if err != nil {
			return
		}
	}

}
func (s *Story) AddPurchase(p Purchase) {
	p.Id = s.LastId + 1
	p.TookAt = time.Now()

	s.Purchases = append(s.Purchases, p)
	s.LastId = p.Id

	s.Save()
}
func (s *Story) DelPurchase(id int) {
	for i, pur := range s.Purchases {
		if pur.Id == id {
			if i == len(s.Purchases)-1 {
				s.Purchases = s.Purchases[:i]
			} else {
				s.Purchases = append(s.Purchases[:i], s.Purchases[i+1:]...)
			}
			s.Save()
		}
	}
}
func (s *Story) EndPurchase(id int) {
	for _, pur := range s.Purchases {
		if pur.Id == id {
			pur.EndAt = time.Now()
			s.Save()
			return
		}
	}

}
func (s *Story) GetAll() []byte {
	return utils.MarshalThis(s.Purchases)
}
func (s *Story) GetByBook(id int) []byte {
	var res []Purchase
	for _, pur := range s.Purchases {
		if pur.BookId == id {
			res = append(res, pur)
		}
	}
	return utils.MarshalThis(res)
}
func (s *Story) GetById(id int) []byte {
	for _, pur := range s.Purchases {
		if pur.Id == id {
			return utils.MarshalThis(pur)
		}
	}
	return nil
}
func (s *Story) GetByUser(id int) []byte {
	var res []Purchase
	for _, pur := range s.Purchases {
		if pur.UserId == id {
			res = append(res, pur)
		}
	}
	return utils.MarshalThis(res)
}
func (s *Story) Save() {
	data, err := json.Marshal(s)
	if err != nil {
		return
	}

	if err := os.WriteFile("./storage/purchases.json", data, 0644); err != nil {
		return
	}
}
func (s *Story) UpdatePurchase(p Purchase) {
	for i, pur := range s.Purchases {
		if pur.Id == p.Id {
			s.Purchases[i] = p
			s.Save()
			return
		}
	}
}
