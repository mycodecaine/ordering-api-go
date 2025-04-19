package server

import (
	"ORDERING-API/internal/bootstrap"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ORDERING-API/docs"
)

func SetupRouter(app *bootstrap.AppContainer) *gin.Engine {
	r := gin.Default()

	// Auth-protected routes
	authorized := r.Group("/")
	authorized.Use(app.AuthMiddleware.MiddlewareFunc())
	{
		authorized.POST("/orders", app.OrderController.CreateOrder)
		authorized.PUT("/orders", app.OrderController.UpdateOrder)
		authorized.GET("/orders", app.OrderController.GetOrder)
	}

	// Auth token
	r.POST("/auth/token", app.AuthController.GetToken)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
