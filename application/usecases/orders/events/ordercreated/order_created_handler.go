package eventhandlers

// application/eventhandlers/order_created_handler.go

import (
	"ORDERING-API/domain/events"
	"log"
)

type OrderCreatedHandler struct{}

func (h OrderCreatedHandler) Handle(event events.DomainEvent) {
	e, ok := event.(events.OrderCreatedEvent)
	if !ok {
		return
	}

	log.Printf("Handling OrderCreatedEvent: OrderID=%s, Timestamp=%v", e.OrderID, e.Timestamp)
}
