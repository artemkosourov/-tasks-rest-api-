package main

import (
	"fmt"
	"log"
	"net/http"
	"tasks-rest-api/config"
	"tasks-rest-api/internal/client"
	"tasks-rest-api/internal/kafka"
	"tasks-rest-api/internal/repository"
	"tasks-rest-api/internal/server"
	"tasks-rest-api/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	repo, err := repository.NewPostgres(cfg)
	if err != nil {
		log.Fatal("Could not connect to DB", err)
	}

	usersClient := client.NewUsersAPIClient("https://jsonplaceholder.typicode.com") // рофло данные
	kafkaProducer := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	defer kafkaProducer.Close()

	svc := service.NewTaskService(repo, kafkaProducer, usersClient)

	mux := http.NewServeMux()
	server.SetupRoutes(mux, svc)

	fmt.Printf("Server is running on %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, mux))
}
