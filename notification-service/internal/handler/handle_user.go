package handler

import (
	"encoding/json"
	"log"

	"github.com/Claudio712005/go-microservices-architecture/notification-service/internal/domain"
	"github.com/Claudio712005/go-microservices-architecture/notification-service/internal/messages"
	"github.com/streadway/amqp"
)

func HandleUserCreated(msg amqp.Delivery) {
	log.Printf("Mensagem recebida em 'user.created': %s", msg.Body)
	log.Println("Enviando email de boas-vindas...")

	var usuario domain.Usuario

	if err := json.Unmarshal(msg.Body, &usuario); err != nil {
		log.Printf("Erro ao deserializar mensagem: %v", err)
	}

	log.Printf("Enviando email de boas-vindas para: %s", usuario.Email)

	if err := messages.SendWelcomeEmailMessage(usuario.Email); err != nil {
		log.Printf("Erro ao enviar email de boas-vindas: %v", err)
	}

	msg.Ack(false)
}
