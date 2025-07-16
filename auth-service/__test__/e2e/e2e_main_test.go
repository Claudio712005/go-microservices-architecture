package e2e

import (
	"os"
	"testing"

	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/config"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/router"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func setUpRoutes() {
    gin.SetMode(gin.ReleaseMode)

    r = gin.New()    
    r.Use(gin.Recovery()) 

    router.SetupRoutes(r, nil)
}

func TestMain(m *testing.M) {
    os.Setenv("APP_ENV", "test")

    config.Load()

    setUpRoutes()

    code := m.Run()

    config.ResetTestDB()
    os.Exit(code)
}
