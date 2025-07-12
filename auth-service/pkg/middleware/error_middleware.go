package middleware

import (
	"errors"
	"fmt"

	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/error"
	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context){
		c.Next()

		if len(c.Errors) == 0 {
			return 
		}

		err := c.Errors.Last().Err
		var appErr *error.AppError

		if errors.As(err, &appErr) {
			c.AbortWithStatusJSON(appErr.Status, appErr)
			return
		}

		appErr = error.Internal("INTERNAL_ERROR", "erro inesperado", err)
		c.AbortWithStatusJSON(appErr.Status, appErr)
	}
}

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, rec any) {
		appErr := error.Internal("PANIC", "erro inesperado", fmt.Errorf("%v", rec))
		c.AbortWithStatusJSON(appErr.Status, appErr)
	})
}