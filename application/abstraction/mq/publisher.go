package mq

type MessageQueuePublisher interface {
	Publish(topic string, message []byte) error
	Close() error
}
