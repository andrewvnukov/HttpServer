package model

import (
	"encoding/json"
	"errors"
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
	Total     int        `json:"total"`
}

type StoryHandler interface {
	Get()
	Save() error
	GetAll() []byte
	GetByUser(int) []byte
	GetByBook(int) []byte
	GetById(int) []byte
	AddPurchase(Purchase) error
	EndPurchase(int) error
	DelPurchase(int) error
	DelPurchaseByBook(int) error
	DelPurchaseByUser(int) error
	UpdatePurchase(Purchase) error
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
func (s *Story) AddPurchase(p Purchase) error {
	p.Id = s.Total
	p.TookAt = time.Now()

	s.Purchases = append(s.Purchases, p)
	s.Total++

	return s.Save()
}
func (s *Story) DelPurchase(id int) error {
	for i, pur := range s.Purchases {
		if pur.Id == id {
			if i == len(s.Purchases)-1 {
				s.Purchases = s.Purchases[:i]
			} else {
				s.Purchases = append(s.Purchases[:i], s.Purchases[i:]...)
				s.Total--
			}
			s.CheckStory()
			err := s.Save()
			if err != nil {
				return err
			} else {
				return nil
			}
		}
	}
	return errors.New("there is no that id")
}
func (s *Story) DelPurchaseByBook(id int) error {
	temp := []Purchase{}
	for _, pur := range s.Purchases {
		if pur.BookId != id {
			temp = append(temp, pur)

		}
	}
	s.Purchases = temp
	s.Total = len(s.Purchases)
	s.CheckStory()
	err := s.Save()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Story) DelPurchaseByUser(id int) error {
	temp := []Purchase{}
	for _, pur := range s.Purchases {
		if pur.UserId != id {
			temp = append(temp, pur)
		}
	}
	s.Purchases = temp
	s.Total = len(s.Purchases)
	s.CheckStory()
	err := s.Save()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (s *Story) EndPurchase(id int) error {
	for i, pur := range s.Purchases {
		if pur.Id == id {
			s.Purchases[i].EndAt = time.Now()
			err := s.Save()
			if err != nil {
				return err
			} else {
				return nil
			}

		}
	}
	return errors.New("there is no id")
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
func (s *Story) Save() error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	if err := os.WriteFile("./storage/purchases.json", data, 0644); err != nil {
		return err
	}

	return nil
}
func (s *Story) UpdatePurchase(p Purchase) error {
	for i, pur := range s.Purchases {
		if pur.Id == p.Id {
			s.Purchases[i] = p
			err := s.Save()
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("not found")
}

func (s *Story) CheckStory() {
	for i := range s.Purchases {
		s.Purchases[i].Id = i
	}
	s.Total = len(s.Purchases)
}
