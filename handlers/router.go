package handlers

import (
	"github.com/gorilla/mux"
)

func BuildRouter(configHandler ConfigHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/configs", configHandler.Create).Methods("POST")
	r.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")
	r.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE")
	return r
}


