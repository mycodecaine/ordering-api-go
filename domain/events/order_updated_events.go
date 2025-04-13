package events

import "time"

type OrderUpdatedEvent struct {
	OrderID   string
	Timestamp time.Time
}

func (e OrderUpdatedEvent) EventType() string {
	return "OrderUpdated"
}

func (e OrderUpdatedEvent) OccurredOn() time.Time {
	return e.Timestamp
}
