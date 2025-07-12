package main

import (
	"log"

	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/config"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/Claudio712005/go-microservices-architecture/auth-service/docs"
)

// @title Auth Service API
// @version 1.0
// @description Este é o serviço de autenticação do sistema de microserviços.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente padrão")
	}

	config.ConectarBanco(false)

	route := gin.Default()

	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))


	router.SetupRoutes(route)

	route.Run(":8080")
}
