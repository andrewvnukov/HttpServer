package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"restapi/model"
	"strconv"
)

type PurchaseHandler struct {
	Purchase model.StoryHandler
}

func NewPurchaseHandler() http.Handler {
	var p PurchaseHandler
	p.Purchase = model.StoryInit()

	return &p
}

func (h *PurchaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if r.URL.Path == "add" {
			m := model.Purchase{}
			if r.Header.Get("Content-Type") == "application/json" {

				data, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				err = json.Unmarshal(data, &m)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					h.Purchase.AddPurchase(m)
				}
			} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
				bookId, err := strconv.Atoi(r.FormValue("book_id"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				userId, err := strconv.Atoi(r.FormValue("user_id"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				m := model.Purchase{
					BookId: bookId,
					UserId: userId,
				}
				h.Purchase.AddPurchase(m)
			}
		} else {
			w.Write([]byte("Not found"))
		}
	case "GET":
		switch r.URL.Path {
		case "getbyid":

			purId, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Add("Content-Type", "Application-Json")
			w.Write(h.Purchase.GetById(purId))
		case "getbybook":
			bookId, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Add("Content-Type", "Application-Json")
			w.Write(h.Purchase.GetByBook(bookId))
		case "getbyuser":
			userId, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Add("Content-Type", "Application-Json")
			w.Write(h.Purchase.GetByBook(userId))
		case "getall":
			w.Header().Add("Content-Type", "Application-Json")
			w.Write(h.Purchase.GetAll())
		default:
			w.Write([]byte("Not found"))
		}
	case "PUT":
		switch r.URL.Path {
		case "update":
			if r.Header.Get("Content-Type") == "application/json" {
				m := model.Purchase{}
				data, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if err = json.Unmarshal(data, &m); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					h.Purchase.UpdatePurchase(m)
				}
			} else if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
				id, err := strconv.Atoi(r.FormValue("id"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				bookId, err := strconv.Atoi(r.FormValue("book_id"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				userId, err := strconv.Atoi(r.FormValue("user_id"))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					m := model.Purchase{
						Id:     id,
						BookId: bookId,
						UserId: userId,
					}
					h.Purchase.UpdatePurchase(m)
				}
			}
		case "endpurchase":
			id, err := strconv.Atoi(r.FormValue("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				h.Purchase.EndPurchase(id)
			}
		default:
			w.Write([]byte("Not found"))
		}
	case "DELETE":
		if r.URL.Path == "delete" {
			id, err := strconv.Atoi(r.FormValue("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			h.Purchase.DelPurchase(id)
		} else {
			w.Write([]byte("Not found"))
		}
	}
}
