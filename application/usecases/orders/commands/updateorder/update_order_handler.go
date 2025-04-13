package commands

import (
	"ORDERING-API/domain/entities"
	"ORDERING-API/domain/events"
	"ORDERING-API/domain/repositories"
)

type UpdateOrderHandler struct {
	repo            repositories.OrderRepository
	eventDispatcher events.EventDispatcher
}

func NewUpdateOrderHandler(repo repositories.OrderRepository, eventDispatcher events.EventDispatcher) *UpdateOrderHandler {
	return &UpdateOrderHandler{repo: repo, eventDispatcher: eventDispatcher}
}

func (h *UpdateOrderHandler) Handle(cmd UpdateOrderCommand) (*UpdateOrderResponse, error) {
	var orderItems []entities.OrderItem

	for _, p := range cmd.OrderItems {
		orderItems = append(orderItems, *entities.NewOrderItem(p.ProductID, p.Quantity)) // Directly append the value
	}

	updateOrder := entities.UpdateOrder(cmd.Id, orderItems, cmd.Notes, cmd.Total)
	err := h.repo.UpdateOrder(updateOrder)

	if err != nil {
		return nil, err
	}

	// Step 2: Dispatch domain events
	h.eventDispatcher.Dispatch(updateOrder.GetEvents())

	// Step 3: Clear them from the aggregate
	updateOrder.ClearEvents()

	return &UpdateOrderResponse{}, nil
}
