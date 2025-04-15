// application/mq/message_handler.go
package mq

type MessageHandler interface {
	Handle(msg []byte) error
}
