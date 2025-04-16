package commands

type UpdateOrderCommand struct {
	Id         string
	OrderItems []OrderItemUpdateDTO
	Notes      string
	Total      float64
	UpdatedBy  string
}

type OrderItemUpdateDTO struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
