package valize

import (
	"github.com/streadway/amqp"
)

type RabbitMQStraegy struct {
	name string
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQStrategy(connURI string, queueName string) *RabbitMQStraegy {
	client, err := amqp.Dial(connURI)
	LogOnErr(err, "Failed to connect to RabbitMQ")
	channel, chanErr := client.Channel()
	LogOnErr(chanErr, "Failed to open a channel")
	_, queueErr := buildQueue(channel, queueName)
	LogOnErr(queueErr, "Failed to declare queue")

	return &RabbitMQStraegy{
		name: queueName,
		conn: client,
		ch:   channel,
	}
}

func (s *RabbitMQStraegy) Push(elem []byte) error {
	publishErr := s.ch.Publish(
		"",     // exchange
		s.name, // routing key
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
	delivery, ok, err := s.ch.Get(
		s.name, // queue
		false,  // auto-ack
	)
	if ok {
		// Requeue message
		delivery.Reject(true)
		return delivery.Body, err
	} else {
		LogOnErr(err, "Failed to register a consumer")
		return nil, err
	}
}

func (s *RabbitMQStraegy) Pop() ([]byte, error) {
	delivery, ok, err := s.ch.Get(
		s.name, // queue
		true,   // auto-ack
	)
	if ok {
		return delivery.Body, err
	} else {
		LogOnErr(err, "Failed to register a consumer")
		return nil, err
	}
}

func (s *RabbitMQStraegy) Clear() error {
	s.ch.QueueDelete(
		s.name,
		true,  // ifused
		false, // ifempty
		true,  // noWait
	)
	// Re-create queue
	_, queueErr := buildQueue(s.ch, s.name)
	LogOnErr(queueErr, "Failed to create queue")
	return queueErr
}

func (s *RabbitMQStraegy) Close() error {
	chanCloseErr := s.ch.Cancel("", true)
	LogOnErr(chanCloseErr, "Failed to close a channel")

	return nil
}

func buildQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	q, queueErr := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return q, queueErr
}
