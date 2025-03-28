package queries

type GetOrderByIdResponse struct {
	Id         string            `json:"id"`
	OrderItems []OrderItemGetDTO `json:"orderItems"`
	Notes      string            `json:"notes"`
	Total      float64           `json:"total"`
}

type OrderItemGetDTO struct {
	Id        string `json:"id"`
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
