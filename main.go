package main

import (
	createorder "ORDERING-API/application/usecases/orders/commands/createorder"
	updateorder "ORDERING-API/application/usecases/orders/commands/updateorder"
	getorderbyid "ORDERING-API/application/usecases/orders/queries/getorderbyid"
	"ORDERING-API/infrastructure/persistence"
	"ORDERING-API/presentation/controllers"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	// Import the generated Swagger docs
	_ "ORDERING-API/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load database connection from environment variables
	connStr := "host=localhost port=5432 dbname=orderingDB user=doadmin password=ipeadmin123456 sslmode=disable"

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err) // Proper error logging
	}
	defer db.Close()

	// Verify the connection is actually established
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database!")

	// Initialize repository
	orderRepo := persistence.NewOrderRepository(db)

	// Initialize handlers
	createOrderHandler := createorder.NewCreateOrderHandler(orderRepo)
	getOrderHandler := getorderbyid.NewGetOrderHandler(orderRepo)
	updateOrderHandler := updateorder.NewUpdateOrderHandler(orderRepo)

	// Initialize controller
	orderController := controllers.NewOrderController(createOrderHandler, getOrderHandler, updateOrderHandler)

	// Initialize Gin router
	r := gin.Default()

	// API routes
	r.POST("/orders", orderController.CreateOrder)
	r.PUT("/orders", orderController.UpdateOrder)
	r.GET("/orders", orderController.GetOrder)

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Run the server on port 8080
	r.Run(":8080")
	// Start server
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
