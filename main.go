package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"projekat/handlers"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"
)

func main() {
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)
    groupRepo := repositories.NewConfigGroupInMemRepository()
    groupService := services.NewConfigGroupService(groupRepo)
    config := model.NewConfig("db_config", 2)
    config.AddParameter("username", "pera")
    config.AddParameter("port", "5432")
    _ = service.Add(config)
    handler := handlers.NewConfigHandler(service)
    groupHandler := handlers.NewConfigGroupHandler(groupService)

    router := handlers.BuildRouter(handler, groupHandler)

	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*1e9)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
