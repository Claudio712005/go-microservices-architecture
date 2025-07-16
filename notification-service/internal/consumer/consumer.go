package consumer

import (
	"log"

	"github.com/Claudio712005/go-microservices-architecture/notification-service/internal/handler"
)

// StartConsumer inicia o consumidor RabbitMQ e aguarda mensagens
func StartConsumer() error {
	conn, ch, err := ConnectRabbitMQ()
	if err != nil {
		return err
	}
	log.Println("Conectado ao RabbitMQ com sucesso")

	defer conn.Close()
	defer ch.Close()

	msgs, err := SetupQueue(ch, "user.events", "user.welcome.email", "user.created")
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler.HandleUserCreated(msg)
		}
	}()

	log.Println("[*] Aguardando mensagens...")

	select {}
}
