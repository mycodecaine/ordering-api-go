package commands

import (
	"ORDERING-API/domain/entities"
	"ORDERING-API/domain/repositories"
)

type UpdateOrderHandler struct {
	repo repositories.IOrderRepository
}

func NewUpdateOrderHandler(repo repositories.IOrderRepository) *UpdateOrderHandler {
	return &UpdateOrderHandler{repo: repo}
}

func (h *UpdateOrderHandler) Handle(cmd UpdateOrderCommand) (*UpdateOrderResponse, error) {
	var orderItems []entities.OrderItem

	for _, p := range cmd.OrderItems {
		orderItems = append(orderItems, *entities.NewOrderItem(p.ProductID, p.Quantity)) // Directly append the value
	}

	newOrder := entities.UpdateOrder(cmd.Id, orderItems, cmd.Notes, cmd.Total)
	err := h.repo.UpdateOrder(newOrder)

	if err != nil {
		return nil, err
	}
	return &UpdateOrderResponse{}, nil
}
