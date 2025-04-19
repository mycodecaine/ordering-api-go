package main

import (
	"ORDERING-API/internal/bootstrap"
	"ORDERING-API/internal/server"
	"log"
)

func main() {
	// Initialize everything
	app := bootstrap.InitializeApp()

	// Start MQ consumer in background
	go func() {
		if err := app.MQConsumer.Consume(); err != nil {
			log.Fatalf("Failed to start MQ consumer: %v", err)
		}
	}()

	// Start HTTP server
	r := server.SetupRouter(app)
	log.Println("Server running at :8080")
	r.Run(":8080")
}
