package main

import (
	"log"

	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/config"
)

func main() {
    log.Println("Starting user-worker...")
    config.Load()
    log.Println("user-worker has been started.")
}