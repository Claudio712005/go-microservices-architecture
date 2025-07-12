package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Envelope padronizado
type Envelope[T any] struct {
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

// Ok representa uma resposta de sucesso genérica
func OK[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Envelope[T]{Data: data})
}

// Created representa uma resposta de criação bem-sucedida
func Created[T any](c *gin.Context, location string, data T) {
	c.Header("Location", location)
	c.JSON(http.StatusCreated, Envelope[T]{Data: data})
}

// NoContent representa uma resposta sem conteúdo
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}