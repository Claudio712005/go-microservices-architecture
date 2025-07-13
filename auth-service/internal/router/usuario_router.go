package router

import (
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/handler"
)

func getUsuarioRoutes() []Route {

	return []Route{
		{
			Path:        "",
			Method:      "POST",
			HandlerFunc: handler.HandleCadastrarUsuario,
			HasAuth:     false,
		},
		{
			Path: "/login",
			Method:      "POST",
			HandlerFunc: handler.HandleLoginUsuario,
			HasAuth:     false,
		},
	}
}
