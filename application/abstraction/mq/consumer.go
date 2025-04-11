package mq

type MessageQueueConsumer interface {
	Consume(queueName string, handler func([]byte)) error
	Close() error
}
