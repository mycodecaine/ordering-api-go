package integrationeventhandlers

import (
	"ORDERING-API/domain/events"
	"encoding/json"
	"log"
)

type SendWhatsappOnOrderUpdatedConsumerHandler struct{}

func (h SendWhatsappOnOrderUpdatedConsumerHandler) Handle(msg []byte) error {
	var event events.OrderUpdatedEvent
	err := json.Unmarshal(msg, &event)
	if err != nil {
		log.Printf("Failed to unmarshal OrderUpdatedEvent: %v", err)
		return err
	}

	log.Printf("Send Whatsapp Consumed OrderUpdatedEvent: OrderID=%s, Timestamp=%v", event.OrderID, event.Timestamp)
	return nil
}
