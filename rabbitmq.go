package valize

import (
	"github.com/streadway/amqp"
)

type RabbitMQStraegy struct {
	conn *amqp.Connection
	name string
}

func NewRabbitMQStrategy(connURI string, queueName string) *RabbitMQStraegy {
	client, err := amqp.Dial(connURI)
	LogOnErr(err, "Failed to connect to RabbitMQ")

	return &RabbitMQStraegy{
		conn: client,
		name: queueName,
	}
}

func (s *RabbitMQStraegy) Push(elem []byte) error {
	ch, q := s.buildQueue()

	publishErr := ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        elem,
		})
	LogOnErr(publishErr, "Failed to publish a message")
	return publishErr
}
func (s *RabbitMQStraegy) Peek() ([]byte, error) {
	return []byte{}, nil
}

func (s *RabbitMQStraegy) Pop() ([]byte, error) {
	ch, q := s.buildQueue()

	// msgs, err := ch.Consume(
	// 	q.Name, // queue
	// 	"",     // consumer
	// 	true,   // auto-ack
	// 	false,  // exclusive
	// 	false,  // no-local
	// 	true,   // no-wait
	// 	nil,    // args
	// )
	delivery, ok, err := ch.Get(
		q.Name, // queue
		true,   // auto-ack
	)
	LogOnErr(err, "Failed to register a consumer")
	if ok {
		return delivery.Body, err
	} else {
		return nil, err
	}
}

func (s *RabbitMQStraegy) Clear() error {
	ch, chanErr := s.conn.Channel()
	LogOnErr(chanErr, "Failed to open a channel")
	ch.QueueDelete(
		s.name,
		true,  // ifused
		false, // ifempty
		true,  // noWait
	)
	return nil
}

func (s *RabbitMQStraegy) buildQueue() (*amqp.Channel, amqp.Queue) {
	ch, chanErr := s.conn.Channel()
	LogOnErr(chanErr, "Failed to open a channel")

	q, queueErr := ch.QueueDeclare(
		s.name, // name
		false,  // durable
		false,  // delete when usused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	LogOnErr(queueErr, "Failed to declare a queue")
	return ch, q
}
