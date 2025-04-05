package entities

import (
	"ORDERING-API/domain/aggregate"
	"ORDERING-API/domain/events"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id         string      `json:"id"`
	OrderItems []OrderItem `json:"orderItems"`
	Notes      string      `json:"notes"`
	Total      float64     `json:"total"`
	aggregate.AggregateRoot
}

func NewOrder(items []OrderItem, notes string, total float64) *Order {
	id := uuid.New().String()
	order := &Order{Id: id, OrderItems: items, Notes: notes, Total: total}

	order.RecordEvent(events.OrderCreatedEvent{
		OrderID:   id,
		Timestamp: time.Now(),
	})
	return order

}

func UpdateOrder(id string, items []OrderItem, notes string, total float64) *Order {
	return &Order{Id: id, OrderItems: items, Notes: notes, Total: total}
}

// Read-only methods
func (o *Order) ID() string { return o.Id }
