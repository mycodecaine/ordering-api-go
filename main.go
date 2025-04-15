// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your token.
package main

import (
	eventhandlers "ORDERING-API/application/events"
	createorder "ORDERING-API/application/usecases/orders/commands/createorder"
	updateorder "ORDERING-API/application/usecases/orders/commands/updateorder"

	integrationordercreatedeventhandlers "ORDERING-API/application/usecases/orders/integrationevents/ordercreated"
	integrationorderupdatedeventhandlers "ORDERING-API/application/usecases/orders/integrationevents/orderupdated"
	getorderbyid "ORDERING-API/application/usecases/orders/queries/getorderbyid"
	"ORDERING-API/infrastructure/auth"
	"ORDERING-API/infrastructure/eventdispatcher"
	"ORDERING-API/infrastructure/mq"
	"ORDERING-API/infrastructure/persistence"
	"ORDERING-API/presentation/controllers"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// Import the generated Swagger docs
	_ "ORDERING-API/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	log.SetOutput(os.Stdout)
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

	// initialize dispatcher
	dispatcher := eventdispatcher.NewSimpleDispatcher()

	publisher, err := mq.NewRabbitMQPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer publisher.Close()

	eventhandler := eventhandlers.NewEventHandler(publisher)

	// Register event handlers
	dispatcher.Register("OrderCreated", eventhandler)
	dispatcher.Register("OrderUpdated", eventhandler)

	// Setup RabbitMQ consumer
	consumer, err := mq.NewRabbitMQConsumer("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer consumer.Close()

	sendemailintegrationeventhandler := integrationordercreatedeventhandlers.SendEmailOnOrderCreatedConsumerHandler{}
	sendwhatsappintegrationeventhandler := integrationordercreatedeventhandlers.SendWhatsappOnOrderCreatedConsumerHandler{}
	sendwhatsapponupdatedintegrationeventhandler := integrationorderupdatedeventhandlers.SendWhatsappOnOrderUpdatedConsumerHandler{}

	consumer.RegisterHandler("OrderCreated", sendemailintegrationeventhandler)
	consumer.RegisterHandler("OrderCreated", sendwhatsappintegrationeventhandler)
	consumer.RegisterHandler("OrderUpdated", sendwhatsapponupdatedintegrationeventhandler)

	// Start consuming
	// Run consumer in a goroutine
	go func() {
		if err := consumer.Consume(); err != nil {
			log.Fatalf("Failed to start consumer: %v", err)
		}
	}()

	// Initialize repository
	orderRepo := persistence.NewOrderRepository(db)

	// Initialize handlers
	createOrderHandler := createorder.NewCreateOrderHandler(orderRepo, dispatcher)
	getOrderHandler := getorderbyid.NewGetOrderHandler(orderRepo)
	updateOrderHandler := updateorder.NewUpdateOrderHandler(orderRepo, dispatcher)

	// Initialize controller
	orderController := controllers.NewOrderController(createOrderHandler, getOrderHandler, updateOrderHandler)

	// Keycloak
	authController := controllers.NewAuthController("http://localhost:7080")

	keycloakMiddleware, err := auth.NewKeycloakMiddleware("http://localhost:7080/realms/agogo", "agogo-client")
	if err != nil {
		log.Fatalf("Failed to initialize Keycloak middleware: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// API routes
	//r.POST("/orders", orderController.CreateOrder)
	//r.PUT("/orders", orderController.UpdateOrder)
	//r.GET("/orders", orderController.GetOrder)

	authorized := r.Group("/")
	authorized.Use(keycloakMiddleware.MiddlewareFunc())
	{
		authorized.POST("/orders", orderController.CreateOrder)
		authorized.PUT("/orders", orderController.UpdateOrder)
		authorized.GET("/orders", orderController.GetOrder)
	}
	r.POST("/auth/token", authController.GetToken)
	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Run the server on port 8080
	r.Run(":8080")
	// Start server
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
