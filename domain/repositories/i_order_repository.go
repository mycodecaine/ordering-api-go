package repositories

import "ORDERING-API/domain/entities"

type IOrderRepository interface {
	SaveOrder(order *entities.Order) (string, error)
	GetOrderByID(id string) (*entities.Order, error)
	DeleteOrder(id string) error
	UpdateOrder(order *entities.Order) error
	GetOrdersWithPagination(limit, offset int) ([]entities.Order, error)
}
