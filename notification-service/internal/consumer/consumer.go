package consumer

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

// StartConsumer é a função que inicia o consumidor de mensagens do RabbitMQ.
func StartConsumer() error {
	amqpUser := os.Getenv("AMQP_USER")
	amqpPassword := os.Getenv("AMQP_PASSWORD")
	amqpURL := os.Getenv("AMQP_URL")

	conn, err := amqp.Dial("amqp://" + amqpUser + ":" + amqpPassword + "@" + amqpURL + "/")
	if err != nil {
		return err
	}
	log.Println("Connected to RabbitMQ successfully")

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	exchange := "user.events"
	err = ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}

	queue := "user.welcome.email"
	q, err := ch.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	routingKey := "user.created"
	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			log.Printf("[x] Mensagem recebida: %s", d.Body)

			log.Printf("[x] Enviando email de boas-vindas para o usuário")

			d.Ack(false)
		}
	}()

	select {}
}