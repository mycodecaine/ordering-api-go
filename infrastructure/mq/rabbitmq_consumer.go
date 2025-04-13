package mq

import (
	"ORDERING-API/application/abstraction/mq"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn          *amqp.Connection
	channel       *amqp.Channel
	queueHandlers map[string][]mq.MessageHandler
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
		conn:          conn,
		channel:       ch,
		queueHandlers: make(map[string][]mq.MessageHandler),
	}, nil
}

func (c *RabbitMQConsumer) Consume() error {
	for queueName, handlers := range c.queueHandlers {

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
			q.Name, "", true, false, false, false, nil,
		)
		if err != nil {
			return err
		}

		// Each queue gets its own goroutine
		go func(queue string, msgs <-chan amqp.Delivery, handlers []mq.MessageHandler) {
			log.Printf("Consuming from queue: %s", queue)
			for msg := range msgs {
				for _, handler := range handlers {
					go func(h mq.MessageHandler, m []byte) {
						if err := h.Handle(m); err != nil {
							log.Printf("Handler error on %s: %v", queue, err)
						}
					}(handler, msg.Body)
				}
			}
		}(queueName, msgs, handlers)
	}

	select {} // block forever
}

// Register handlers
func (c *RabbitMQConsumer) RegisterHandler(queueName string, handler mq.MessageHandler) {
	c.queueHandlers[queueName] = append(c.queueHandlers[queueName], handler)
}

func (r *RabbitMQConsumer) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}
	return r.conn.Close()
}
