package events

import "time"

type OrderCreatedEvent struct {
	OrderID   string
	Timestamp time.Time
}

func (e OrderCreatedEvent) EventType() string {
	return "OrderCreated"
}

func (e OrderCreatedEvent) OccurredOn() time.Time {
	return e.Timestamp
}
