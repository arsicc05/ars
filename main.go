package main

import (
	"net/http"
	"projekat/handlers"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"

	"github.com/gorilla/mux"
)

func main() {
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)
    config := model.NewConfig("db_config", 2)
    config.AddParameter("username", "pera")
    config.AddParameter("port", "5432")
    _ = service.Add(config)
	handler := handlers.NewConfigHandler(service)

	router := mux.NewRouter()

	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")

	http.ListenAndServe("0.0.0.0:8000", router)
}
