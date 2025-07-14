package config

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func AmqpConfig() {
	amqpUser := os.Getenv("AMQP_USER")
	amqpPassword := os.Getenv("AMQP_PASSWORD")
	amqpURL := os.Getenv("AMQP_URL")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", amqpUser, amqpPassword, amqpURL))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		panic(err)
	}
	
	log.Println("Connected to RabbitMQ successfully")

	defer conn.Close()
}