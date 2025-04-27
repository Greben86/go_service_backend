package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// API приложения.
type API struct {
	r *mux.Router // маршрутизатор запросов
	// db *db.DB     // база данных
}

// Конструктор API.
func ApiNewInstance() *API {
	api := API{}
	// api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

func (api *API) exapmpleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("123")
	w.Write([]byte("example"))
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.HandleFunc("/example", api.exapmpleHandler).Methods(http.MethodGet)
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}
