package valize

import (
	"errors"
	"log"

	"github.com/streadway/amqp"
)

type RabbitMQStraegy struct {
	conn  *amqp.Connection
	name  string
	queue amqp.Queue
}

func NewRabbitMQStrategy(connURI string, queueName string) *RabbitMQStraegy {
	client, err := amqp.Dial(connURI)
	LogOnErr(err, "Failed to connect to RabbitMQ")
	ch, chanErr := client.Channel()
	LogOnErr(chanErr, "Failed to open a channel")
	q, queueErr := buildQueue(ch, queueName)
	LogOnErr(queueErr, "Failed to declare queue")

	return &RabbitMQStraegy{
		conn:  client,
		name:  queueName,
		queue: q,
	}
}

func (s *RabbitMQStraegy) Push(elem []byte) error {
	ch, chanErr := s.conn.Channel()
	LogOnErr(chanErr, "Failed to open a channel")
	publishErr := ch.Publish(
		"",           // exchange
		s.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
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
	ch, chanErr := s.conn.Channel()
	LogOnErr(chanErr, "Failed to open a channel")
	// msgs, err := s.channel.Consume(
	// 	s.queue.Name, // queue
	// 	"",           // consumer
	// 	true,         // auto-ack
	// 	false,        // exclusive
	// 	false,        // no-local
	// 	true,         // no-wait
	// 	nil,          // args
	// )
	delivery, ok, err := ch.Get(
		s.queue.Name, // queue
		true,         // auto-ack
	)
	if ok {
		log.Printf("%#v delivery", delivery)
		return delivery.Body, err
	} else if err != nil {
		LogOnErr(err, "Failed to register a consumer")
		return nil, err
	} else {
		return nil, errors.New("Queue is empty")
	}
	// delivery := <-msgs
	// return (<-msgs).Body, err
	// return []byte{}, nil
}

func (s *RabbitMQStraegy) Clear() error {
	return nil
}

func (s *RabbitMQStraegy) Close() error {
	log.Println("Closing channel")
	ch, chanErr := s.conn.Channel()
	LogOnErr(chanErr, "Failed to open a channel")
	defer ch.Close()
	ch.QueueDelete(
		s.name,
		true,  // ifused
		false, // ifempty
		true,  // noWait
	)
	chanCloseErr := ch.Cancel("", true)
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
