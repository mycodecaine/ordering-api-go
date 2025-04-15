package integrationeventhandlers

import (
	"ORDERING-API/domain/events"
	"encoding/json"
	"log"
)

type SendEmailOnOrderCreatedConsumerHandler struct{}

func (h SendEmailOnOrderCreatedConsumerHandler) Handle(msg []byte) error {
	var event events.OrderCreatedEvent
	err := json.Unmarshal(msg, &event)
	if err != nil {
		log.Printf("Failed to unmarshal OrderCreatedEvent: %v", err)
		return err
	}

	log.Printf("Send Email Consumed OrderCreatedEvent: OrderID=%s, Timestamp=%v", event.OrderID, event.Timestamp)
	return nil
}
