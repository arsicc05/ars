package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projekat/handlers"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	configRepo := repositories.NewConfigInMemRepository()
	groupRepo := repositories.NewConfigGroupInMemRepository()
	
	configService := services.NewConfigService(configRepo)
	groupService := services.NewConfigGroupService(groupRepo)
	
	config := model.NewConfig("db_config", 2)
	config.AddParameter("username", "pera")
	config.AddParameter("port", "5432")
	config.AddParameter("host", "localhost")
	_ = configService.Add(config)
	
	group := model.NewConfigGroup("web_configs", 1)
	
	webConfig := model.NewGroupConfig("web_server")
	webConfig.AddParameter("port", "8080")
	webConfig.AddParameter("host", "0.0.0.0")
	webConfig.AddLabel("environment", "development")
	webConfig.AddLabel("team", "backend")
	
	group.AddConfig(webConfig)
	_ = groupService.Add(group)
	
	configHandler := handlers.NewConfigHandler(configService)
	groupHandler := handlers.NewConfigGroupHandler(groupService)
	
	router := mux.NewRouter()
	
	router.HandleFunc("/configs", configHandler.GetAll).Methods("GET")
	router.HandleFunc("/configs", configHandler.Create).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Get).Methods("GET")
	router.HandleFunc("/configs/{name}/{version}", configHandler.Delete).Methods("DELETE")
	
	router.HandleFunc("/groups", groupHandler.GetAll).Methods("GET")
	router.HandleFunc("/groups", groupHandler.Create).Methods("POST")
	router.HandleFunc("/groups/{name}/{version}", groupHandler.Get).Methods("GET")
	router.HandleFunc("/groups/{name}/{version}", groupHandler.Delete).Methods("DELETE")
	
	router.HandleFunc("/groups/{name}/{version}/configs", groupHandler.AddConfig).Methods("POST")
	router.HandleFunc("/groups/{name}/{version}/configs/{configName}", groupHandler.GetConfig).Methods("GET")
	router.HandleFunc("/groups/{name}/{version}/configs/{configName}", groupHandler.RemoveConfig).Methods("DELETE")

	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on :8000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on :8000: %v\n", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
