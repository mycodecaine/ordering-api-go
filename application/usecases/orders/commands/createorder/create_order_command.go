package commands

type CreateOrderCommand struct {
	OrderItems []OrderItemCreateDTO
	Notes      string
	Total      float64
}

type OrderItemCreateDTO struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
