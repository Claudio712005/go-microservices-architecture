package config

import (
	"context"
	"fmt"
	"log"

	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/consumer"
	"github.com/joho/godotenv"
)

func Load() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	conectarBanco()

	_, ch, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	consumer := consumer.NewAuditConsumer(ch, DB)

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("consumer stopped with error: %v", err)
		}
	}()
}
