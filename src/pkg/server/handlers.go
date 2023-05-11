package server

import (
	"encoding/json"
	"hezzlService/src/internal/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const ErrNotFound = "errors.NotFound"

func (ws *WebServer) handleDefaultPage(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Home"))
}

func (ws *WebServer) handleItemCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else {
		campID, err := strconv.Atoi(r.URL.Query().Get("campaignId"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("[SERVER] | POST /item/create | can't parse request body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var item models.Item
		item.CampaignID = campID

		err = json.Unmarshal(data, &item)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := ws.createItem(item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Printf("[SERVER] | POST /item/create | create item was succesfully: %v\n", item)
			_, _ = w.Write(res)
		}
	}
}

func (ws *WebServer) handleItemUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else {
		ID, err := strconv.Atoi(r.URL.Query().Get("Id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		campID, err := strconv.Atoi(r.URL.Query().Get("campaignId"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("[SERVER] | PATCH /item/update | can't parse request body: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var item models.Item
		item.ID = ID
		item.CampaignID = campID

		err = json.Unmarshal(data, &item)
		if err != nil || item.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := ws.updateItem(item)
		if err != nil {
			_, _ = w.Write([]byte(ErrNotFound))
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			log.Printf("[SERVER] | PATCH /item/update | succesfull: %v\n", item)
			_, _ = w.Write(res)
		}
	}
}

func (ws *WebServer) handleItemRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else {
		ID, err := strconv.Atoi(r.URL.Query().Get("Id"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		campID, err := strconv.Atoi(r.URL.Query().Get("campaignId"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var data models.Item
		data.ID = ID
		data.CampaignID = campID

		res, err := ws.deleteItem(&data)
		if err != nil {
			if err.Error() == ErrNotFound {
				_, _ = w.Write([]byte(ErrNotFound))
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			log.Println("[SERVER] | DELETE /item/remove | succesfull\n")
			_, _ = w.Write(res)
		}
	}
}

func (ws *WebServer) handleItemsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else {
		res, err := ws.getItems()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			log.Printf("[SERVER] | GET /items/list | succesfull\n")
			_, _ = w.Write(res)
		}
	}
}
