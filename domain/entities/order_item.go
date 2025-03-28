package entities

import "github.com/google/uuid"

type OrderItem struct {
	Id        string `json:"id"`
	OrderID   string `json:"orderId"`
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func NewOrderItem(productId string, quantity int) *OrderItem {
	return &OrderItem{Id: uuid.New().String(), ProductID: productId, Quantity: quantity}
}
