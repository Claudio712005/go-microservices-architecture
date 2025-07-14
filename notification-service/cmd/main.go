package main

import (
	"github.com/Claudio712005/go-microservices-architecture/notification-service/internal/config"
)

func main() {
	config.LoadEnv()

	config.AmqpConfig()

}