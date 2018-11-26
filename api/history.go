package api

import (
	"encoding/json"
	"net/http"
)

func (a Api) priceHistory(w http.ResponseWriter, r *http.Request) {

	siteID := r.URL.Query().Get("site_id")
	productID := r.URL.Query().Get("product_id")

	if siteID == "" {
		http.Error(w, "no siteID parameter provided", http.StatusBadRequest)
		return
	}

	if productID == "" {
		http.Error(w, "no productID parameter provided", http.StatusBadRequest)
		return
	}

	history, err := a.Storage.Get(siteID, productID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(history)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)

}
