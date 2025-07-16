package main

import (
	"log"
	"os"

	_ "github.com/Claudio712005/go-microservices-architecture/auth-service/docs"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/config"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/mq"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/router"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service API
// @version 1.0
// @description Este é o serviço de autenticação do sistema de microserviços.
// @host localhost:8080
// @BasePath /api/v1
func main() {

	config.Load()

	bus, err := mq.NewRabbitBus(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatalf("erro ao conectar RabbitMQ: %v", err)
	}

	defer bus.Close()
	route := gin.Default()

	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.SetupRoutes(route, bus)

	route.Run(":8080")
}
