package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/repository"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/schema"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/service"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

const (
	exchangeName    = "user.events"
	exchangeType    = "topic"
	queueName       = "user.audit"
	routingKeyAudit = "user.audit"
	consumerTag     = ""
)

// AuditConsumer é responsável por consumir eventos de auditoria do RabbitMQ.
type AuditConsumer struct {
	channel *amqp.Channel
	db      *gorm.DB
}

// NewAuditConsumer cria uma nova instância do AuditConsumer.
func NewAuditConsumer(ch *amqp.Channel, db *gorm.DB) *AuditConsumer {
	return &AuditConsumer{channel: ch, db: db}
}

// Start inicia o consumidor de eventos de auditoria.
func (c *AuditConsumer) Start(ctx context.Context) error {
	if err := c.channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if _, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := c.channel.QueueBind(
		queueName,
		routingKeyAudit,
		exchangeName,
		false,
		nil,
	); err != nil {
		return err
	}

	deliveries, err := c.channel.Consume(
		queueName,
		consumerTag,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("✅ AuditConsumer aguardando mensagens em", queueName)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d, ok := <-deliveries:
			if !ok {
				return amqp.ErrClosed
			}
			c.process(d)
		}
	}
}

func (c *AuditConsumer) process(d amqp.Delivery) {
	log.Println("Recebendo evento de auditoria")

	var body schema.UserAuditEvent

	if err := json.Unmarshal(d.Body, &body); err != nil {
		log.Printf("Erro ao deserializar mensagem: %v", err)
		return
	}

	service := service.NewAuditService(repository.NewAuditRepository(c.db))
	if err := service.RegistryNewAudit(&body); err != nil {
		log.Printf("Erro ao registrar evento de auditoria: %v", err)
		return
	}

	log.Printf("Evento de auditoria registrado com sucesso: %s", body.EventType)
}
