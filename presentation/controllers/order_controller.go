package controllers

import (
	"ORDERING-API/application"
	createorder "ORDERING-API/application/usecases/orders/commands/createorder"
	updateorder "ORDERING-API/application/usecases/orders/commands/updateorder"
	getorderbyid "ORDERING-API/application/usecases/orders/queries/getorderbyid"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	createOrderHandler *createorder.CreateOrderHandler
	getOrderHandler    *getorderbyid.GetOrderByIdHandler
	updateOrderHandler *updateorder.UpdateOrderHandler
}

func NewOrderController(createOrder *createorder.CreateOrderHandler, getOrder *getorderbyid.GetOrderByIdHandler, updateOrder *updateorder.UpdateOrderHandler) *OrderController {
	return &OrderController{
		createOrderHandler: createOrder,
		getOrderHandler:    getOrder,
		updateOrderHandler: updateOrder,
	}
}

// CreateOrder creates a new order
// @Summary Create a new order
// @Description Create an order with items and details
// @Tags orders
// @Accept json
// @Produce json
// @Param order body commands.CreateOrderCommand true "Order JSON"
// @Success 201 {object} createorder.CreateOrderResponse
// @Failure 400 {object} application.ErrorResponse
// @Failure 500 {object} application.ErrorResponse
// @Security BearerAuth
// @Router /orders [post]
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var request createorder.CreateOrderCommand
	// Bind JSON request body to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, application.ErrorResponse{Error: err.Error()})
		return
	}
	// Process order creation
	order, err := oc.createOrderHandler.Handle(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, application.ErrorResponse{Error: err.Error()})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, order)
}

// Update order
// @Summary Update order
// @Description Update an order with items and details
// @Tags orders
// @Accept json
// @Produce json
// @Param order body updateorder.UpdateOrderCommand true "Order JSON"
// @Success 200 {object} updateorder.UpdateOrderResponse
// @Failure 400 {object} application.ErrorResponse
// @Failure 500 {object} application.ErrorResponse
// @Security BearerAuth
// @Router /orders [put]
func (oc *OrderController) UpdateOrder(c *gin.Context) {
	var request updateorder.UpdateOrderCommand
	// Bind JSON request body to struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, application.ErrorResponse{Error: err.Error()})
		return
	}
	// Process order update
	order, err := oc.updateOrderHandler.Handle(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, application.ErrorResponse{Error: err.Error()})
		return
	}

	// Success response
	c.JSON(http.StatusOK, order)
}

// GetOrder retrieves an order by ID
// @Summary Get order by ID
// @Description Fetch an order's details
// @Tags orders
// @Produce json
// @Param id query string true "Order ID"
// @Success 200 {object} queries.GetOrderByIdResponse
// @Failure 404 {object} application.ErrorResponse
// @Security BearerAuth
// @Router /orders [get]
func (oc *OrderController) GetOrder(c *gin.Context) {
	orderID := c.Query("id")
	query := getorderbyid.GetOrderByIdQuery{OrderID: orderID}
	order, err := oc.getOrderHandler.Handle(query)
	if err != nil {
		c.JSON(http.StatusNotFound, application.ErrorResponse{Error: err.Error()})
		return
	}
	// Success response
	c.JSON(http.StatusOK, order)
}
