package api

import (
	"encoding/json"
	"github.com/GaruGaru/keeprice/models"
	"net/http"
)

func (a Api) addPrice(w http.ResponseWriter, r *http.Request) {
	productPrice := &models.ProductPrice{}
	err := json.NewDecoder(r.Body).Decode(productPrice)

	if err != nil {
		http.Error(w, "invalid request body "+err.Error(), http.StatusBadRequest)
		return
	}

	err = a.Storage.Store(*productPrice)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte{})
}

