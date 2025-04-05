package commands

import (
	"ORDERING-API/domain/entities"
	"ORDERING-API/domain/events"
	"ORDERING-API/domain/repositories"
)

type CreateOrderHandler struct {
	repo            repositories.IOrderRepository
	eventDispatcher events.EventDispatcher
}

func NewCreateOrderHandler(repo repositories.IOrderRepository, eventDispatcher events.EventDispatcher) *CreateOrderHandler {
	return &CreateOrderHandler{repo: repo, eventDispatcher: eventDispatcher}
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

	// Step 2: Dispatch domain events
	h.eventDispatcher.Dispatch(newOrder.GetEvents())

	// Step 3: Clear them from the aggregate
	newOrder.ClearEvents()

	return &CreateOrderResponse{Id: orderId}, nil
}
