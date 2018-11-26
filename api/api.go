package api

import (
	"github.com/GaruGaru/keeprice/storage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type Api struct {
	Storage storage.PriceStorage
}

var AllowedHeaders = []string{"X-Requested-With", "Accept", "Content-Type"}
var AllowedOrigins = []string{"*"}
var AllowedMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}

func (a Api) initRouting(router *mux.Router) {
	router.HandleFunc("/probe", a.probe)
	router.HandleFunc("/product", a.addPrice).Methods("POST")
	router.HandleFunc("/product", a.priceHistory).Methods("GET")
}

func (a Api) Run(addr string) error {
	router := mux.NewRouter()
	a.initRouting(router)

	allowedHeaders := handlers.AllowedHeaders(AllowedHeaders)
	allowedOrigins := handlers.AllowedOrigins(AllowedOrigins)
	allowedMethods := handlers.AllowedMethods(AllowedMethods)

	return http.ListenAndServe(addr, handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods)(router))
}
