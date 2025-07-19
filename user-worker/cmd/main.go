package main

import (
	"log"

	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/config"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
    log.Println("Starting user-worker...")
    config.Load()
    log.Println("user-worker has been started.")

    route := gin.Default()

    router.SetupRoutes(route)

    route.Run(":8082")

    log.Println("user-worker is running on port 8082")
}