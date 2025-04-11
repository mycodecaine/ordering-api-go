package mq

import (
	"ORDERING-API/application/abstraction/mq"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	handlers []mq.MessageHandler // <- List of handlers
}

func NewRabbitMQConsumer(amqpURL string) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *RabbitMQConsumer) Consume(queueName string) error {
	q, err := c.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			for _, handler := range c.handlers {
				go func(h mq.MessageHandler, m []byte) {
					if err := h.Handle(m); err != nil {
						log.Printf("Handler error: %v", err)
					}
				}(handler, msg.Body)
			}
		}
	}()

	log.Printf("Consumer started on queue: %s", queueName)
	select {} // Block forever

}

// Register handlers
func (c *RabbitMQConsumer) RegisterHandler(handler mq.MessageHandler) {
	c.handlers = append(c.handlers, handler)
}

func (r *RabbitMQConsumer) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}
