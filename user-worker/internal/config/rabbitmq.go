package config

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func connectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	user := os.Getenv("AMQP_USER")
	pass := os.Getenv("AMQP_PASSWORD")
	host := os.Getenv("AMQP_URL")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", user, pass, host))
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}