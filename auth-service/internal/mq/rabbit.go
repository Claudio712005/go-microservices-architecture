package mq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/schema"
	amqp "github.com/streadway/amqp"
)

// RabbitBus é a implementação do EventBus usando RabbitMQ.
type RabbitBus struct {
    conn *amqp.Connection
    ch   *amqp.Channel
}

const exchange = "user.events"
const routingKeyUserCreated = "user.created"

// NewRabbitBus cria uma nova instância do RabbitBus.
func NewRabbitBus(url string) (EventBus, error) {
    conn, err := amqp.Dial(url)
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, err
    }

    if err := ch.ExchangeDeclare(
        exchange, "topic",
        true,  
        false, 
        false, 
        false, 
        nil,   
    ); err != nil {
        conn.Close()
        return nil, err
    }

    return &RabbitBus{conn: conn, ch: ch}, nil
}

// PublishUserCreated publica um evento de usuário criado no RabbitMQ.
func (r *RabbitBus) PublishUserCreated(ctx context.Context, evt schema.UsuarioCreated) error {
    body, _ := json.Marshal(evt)

    headers := amqp.Table{}
    if trace, ok := ctx.Value("trace_id").(string); ok {
        headers["trace_id"] = trace
    }

    return r.ch.Publish(
        exchange,
        routingKeyUserCreated,
        false,
        false,
        amqp.Publishing{
            ContentType:  "application/json",
            Body:         body,
            Timestamp:    time.Now(),
            DeliveryMode: amqp.Persistent,
            Headers:      headers,
        },
    )
}

// Close fecha a conexão e o canal do RabbitMQ.
func (r *RabbitBus) Close() error {
    if err := r.ch.Close(); err != nil {
        r.conn.Close()
        return err
    }
    return r.conn.Close()
}
