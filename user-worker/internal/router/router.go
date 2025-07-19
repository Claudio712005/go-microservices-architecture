package router

import (
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// Route define a estrutura de uma rota
type Route struct {
	Path        string
	Method      string
	HandlerFunc func(*gin.Context)
	HasAuth     bool
}

// SetupRoutes inicializa as rotas da aplicação
func SetupRoutes(r *gin.Engine) {

	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.ErrorMiddleware())

	v1 := r.Group("/api/v1")

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "user-worker is alive"})
	})

	auditGroup := v1.Group("/usuario-audit")
	registerRoutes(auditGroup, getAuditedUsuarioRoutes())
}

func registerRoutes(g *gin.RouterGroup, routes []Route) {
	for _, rt := range routes {

		if rt.HasAuth {
			g.Handle(
				rt.Method,
				rt.Path,
				middleware.AutenticacaoMiddleware(),
				rt.HandlerFunc,
			)
		} else {

			g.Handle(rt.Method, rt.Path, rt.HandlerFunc)
		}
	}
}
