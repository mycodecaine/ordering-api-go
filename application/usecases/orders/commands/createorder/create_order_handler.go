package commands

import (
	"ORDERING-API/domain/entities"
	"ORDERING-API/domain/repositories"
)

type CreateOrderHandler struct {
	repo repositories.IOrderRepository
}

func NewCreateOrderHandler(repo repositories.IOrderRepository) *CreateOrderHandler {
	return &CreateOrderHandler{repo: repo}
}

func (h *CreateOrderHandler) Handle(cmd CreateOrderCommand) (*CreateOrderResponse, error) {
	var orderItems []entities.OrderItem

	for _, p := range cmd.OrderItems {
		orderItems = append(orderItems, *entities.NewOrderItem(p.ProductID, p.Quantity)) // Directly append the value
	}

	newOrder := entities.NewOrder(orderItems, cmd.Notes, cmd.Total)
	orderId, err := h.repo.SaveOrder(newOrder)

	if err != nil {
		return nil, err
	}
	return &CreateOrderResponse{Id: orderId}, nil
}
