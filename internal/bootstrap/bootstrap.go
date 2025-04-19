package bootstrap

import (
	publisher "ORDERING-API/application/abstraction/mq"
	eventhandlers "ORDERING-API/application/events"
	createorder "ORDERING-API/application/usecases/orders/commands/createorder"
	updateorder "ORDERING-API/application/usecases/orders/commands/updateorder"
	getorderbyid "ORDERING-API/application/usecases/orders/queries/getorderbyid"
	"ORDERING-API/domain/repositories"
	"ORDERING-API/infrastructure/auth"
	"ORDERING-API/infrastructure/eventdispatcher"
	"ORDERING-API/infrastructure/mq"
	"ORDERING-API/infrastructure/persistence"
	"ORDERING-API/internal/config"
	"ORDERING-API/presentation/controllers"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"

	integrationordercreatedeventhandlers "ORDERING-API/application/usecases/orders/integrationevents/ordercreated"
	integrationorderupdatedeventhandlers "ORDERING-API/application/usecases/orders/integrationevents/orderupdated"
)

type AppContainer struct {
	DB          *sql.DB
	Dispatcher  *eventdispatcher.InMemoryDispatcher
	OrderRepo   repositories.OrderRepository
	MQPublisher publisher.MessageQueuePublisher
	MQConsumer  *mq.RabbitMQConsumer

	OrderController *controllers.OrderController
	AuthController  *controllers.AuthController
	AuthMiddleware  *auth.KeycloakMiddleware
}

func InitializeApp() *AppContainer {
	initLogger()
	cfg := config.Load()
	// === Database Setup ===
	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword, cfg.DBSSLMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	// === Dispatcher ===
	dispatcher := eventdispatcher.NewSimpleDispatcher()

	// === MQ Setup ===
	publisher, err := mq.NewRabbitMQPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ publisher: %v", err)
	}

	consumer, err := mq.NewRabbitMQConsumer(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ consumer: %v", err)
	}

	// === Register Event Handlers ===
	eventHandler := eventhandlers.NewEventHandler(publisher)
	dispatcher.Register("OrderCreated", eventHandler)
	dispatcher.Register("OrderUpdated", eventHandler)

	consumer.RegisterHandler("OrderCreated", integrationordercreatedeventhandlers.SendEmailOnOrderCreatedConsumerHandler{})
	consumer.RegisterHandler("OrderCreated", integrationordercreatedeventhandlers.SendWhatsappOnOrderCreatedConsumerHandler{})
	consumer.RegisterHandler("OrderUpdated", integrationorderupdatedeventhandlers.SendWhatsappOnOrderUpdatedConsumerHandler{})

	// === Repository & Use Cases ===
	orderRepo := persistence.NewOrderRepository(db)
	createHandler := createorder.NewCreateOrderHandler(orderRepo, dispatcher)
	updateHandler := updateorder.NewUpdateOrderHandler(orderRepo, dispatcher)
	getHandler := getorderbyid.NewGetOrderHandler(orderRepo)

	// === Controllers ===
	orderController := controllers.NewOrderController(createHandler, getHandler, updateHandler)
	authController := controllers.NewAuthController(cfg.KeycloakURL)

	keycloakMiddleware, err := auth.NewKeycloakMiddleware(cfg.KeycloakURL+"/realms/"+cfg.KeycloakRealm, cfg.KeycloakClient)
	if err != nil {
		log.Fatalf("Failed to setup Keycloak middleware: %v", err)
	}

	return &AppContainer{
		DB:              db,
		Dispatcher:      dispatcher,
		OrderRepo:       orderRepo,
		MQPublisher:     publisher,
		MQConsumer:      consumer,
		OrderController: orderController,
		AuthController:  authController,
		AuthMiddleware:  keycloakMiddleware,
	}
}

func initLogger() {
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		DisableQuote:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.DebugLevel)
}
