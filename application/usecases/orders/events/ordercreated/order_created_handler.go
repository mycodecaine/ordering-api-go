package eventhandlers

// application/eventhandlers/order_created_handler.go

import (
	"ORDERING-API/application/abstraction/mq" // your interface
	"ORDERING-API/domain/events"
	"encoding/json"
	"log"
)

type OrderCreatedHandler struct {
	Publisher mq.MessageQueuePublisher
}

func NewOrderCreatedHandler(publisher mq.MessageQueuePublisher) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		Publisher: publisher,
	}
}

func (h OrderCreatedHandler) Handle(event events.DomainEvent) {
	e, ok := event.(events.OrderCreatedEvent)
	if !ok {
		return
	}

	log.Printf("Handling OrderCreatedEvent: OrderID=%s, Timestamp=%v", e.OrderID, e.Timestamp)

	// Marshal the event to JSON
	message, err := json.Marshal(e)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
		return
	}

	// Publish to queue
	err = h.Publisher.Publish("order.created", message)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}
}
