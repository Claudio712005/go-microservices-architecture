package consumer

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

// ConnectRabbitMQ conecta ao RabbitMQ usando as variáveis de ambiente
func ConnectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
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

// SetupQueue configura a fila e o exchange no RabbitMQ
func SetupQueue(ch *amqp.Channel, exchange, queueName, routingKey string) (<-chan amqp.Delivery, error) {
	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if err := ch.QueueBind(q.Name, routingKey, exchange, false, nil); err != nil {
		return nil, err
	}

	return ch.Consume(q.Name, "", false, false, false, false, nil)
}
