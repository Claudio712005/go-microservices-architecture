package main

import (
	"log"

	"github.com/Claudio712005/go-microservices-architecture/notification-service/internal/config"
	"github.com/Claudio712005/go-microservices-architecture/notification-service/internal/consumer"
)

func main() {
	config.LoadEnv()

	err := consumer.StartConsumer()
	if err != nil {
		log.Fatalf("Erro ao iniciar consumidor: %v", err)
	}
}