package eventhandlers

// application/eventhandlers/order_created_handler.go

import (
	"ORDERING-API/application/abstraction/mq" // your interface
	"ORDERING-API/domain/events"
	"encoding/json"
	"log"
)

type EventHandler struct {
	Publisher mq.MessageQueuePublisher
}

func NewEventHandler(publisher mq.MessageQueuePublisher) *EventHandler {
	return &EventHandler{
		Publisher: publisher,
	}
}

func (h EventHandler) Handle(e events.DomainEvent) {
	if e == nil {
		log.Println("Received nil event")
		return
	}

	eventType := e.EventType()
	log.Printf("Handling Event:")

	// Marshal the event to JSON
	message, err := json.Marshal(e)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	// Publish to queue
	err = h.Publisher.Publish(eventType, message)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}
}
