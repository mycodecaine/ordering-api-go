package queries

import (
	"ORDERING-API/domain/repositories"
)

type GetOrderByIdHandler struct {
	repo repositories.OrderRepository
}

func NewGetOrderHandler(repo repositories.OrderRepository) *GetOrderByIdHandler {
	return &GetOrderByIdHandler{repo: repo}
}

func (h *GetOrderByIdHandler) Handle(query GetOrderByIdQuery) (*GetOrderByIdResponse, error) {
	order, err := h.repo.GetOrderByID(query.OrderID)

	if err != nil {
		return nil, err
	}
	var orderItems []OrderItemGetDTO

	for _, p := range order.OrderItems {
		orderItems = append(orderItems, OrderItemGetDTO{Id: p.Id, Quantity: p.Quantity, ProductID: p.ProductID}) // Directly append the value
	}
	orderResponse := GetOrderByIdResponse{Id: order.Id, Notes: order.Notes, Total: order.Total, OrderItems: orderItems}

	return &orderResponse, nil

}
