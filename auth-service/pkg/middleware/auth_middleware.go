package middleware

import (
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/error"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/security"
	"github.com/gin-gonic/gin"
)

// AutenticacaoMiddleware verifica se o token de autenticação é válido
func AutenticacaoMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenBearer := c.GetHeader("Authorization")

		if tokenBearer == "" {
			c.Error(error.Unauthorized("UNAUTHORIZED", "Token de autenticação não fornecido", nil))
			c.Abort()
			return
		}

		token := tokenBearer[len("Bearer "):]

		if err := security.ValidarToken(token); err != nil {
			c.Error(error.Unauthorized("UNAUTHORIZED", "Token de autenticação inválido", err))
			c.Abort()
			return
		}
	}
}