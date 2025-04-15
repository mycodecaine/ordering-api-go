package mq

type MessageQueueConsumer interface {
	Consume(queueName string) error
	Close() error
}
