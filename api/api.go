package api

import (
	"encoding/json"
	"github.com/GaruGaru/keeprice/models"
	"github.com/GaruGaru/keeprice/storage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type Api struct {
	Storage storage.PriceStorage
}

func (a Api) initRouting(router *mux.Router) {
	router.HandleFunc("/probe", a.probe)
	router.HandleFunc("/product", a.addPrice).Methods("POST")
	router.HandleFunc("/product", a.priceHistory).Methods("GET")
}

func (a Api) Run(addr string) error {
	router := mux.NewRouter()
	a.initRouting(router)

	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Accept", "Content-Type"})
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	return http.ListenAndServe(addr, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router))
}

func (a Api) probe(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
}

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
