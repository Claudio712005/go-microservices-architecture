package config

import (
	"context"
	"log"

	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/consumer"
	"github.com/joho/godotenv"
)

func Load() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conectarBanco()

	conn, ch, err := connectRabbitMQ()
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }
    defer conn.Close()
    defer ch.Close()

    consumer := consumer.NewAuditConsumer(ch, DB)
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    if err := consumer.Start(ctx); err != nil {
        log.Fatalf("consumer stopped with error: %v", err)
    }
}