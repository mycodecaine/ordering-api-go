package entities

import (
	"github.com/google/uuid"
)

type Order struct {
	Id         string      `json:"id"`
	OrderItems []OrderItem `json:"orderItems"`
	Notes      string      `json:"notes"`
	Total      float64     `json:"total"`
}

func NewOrder(items []OrderItem, notes string, total float64) *Order {
	id := uuid.New().String()
	return &Order{Id: id, OrderItems: items, Notes: notes, Total: total}
}

func UpdateOrder(id string, items []OrderItem, notes string, total float64) *Order {
	return &Order{Id: id, OrderItems: items, Notes: notes, Total: total}
}

// Read-only methods
func (o *Order) ID() string { return o.Id }
