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
		{
			Path: "/logado",
			Method:      "GET",
			HandlerFunc: handler.HandleBuscarUsuarioLogado,
			HasAuth:     true,
		},
		{
			Path: "/:id",
			Method:      "PUT",
			HandlerFunc: handler.HandleAlterarUsuario,
			HasAuth:     true,
		}, 
		{
			Path: "/senha",
			Method:      "PUT",
			HandlerFunc: handler.HandleAlterarSenhaUsuarioLogado,
			HasAuth:     true,
		},
	}
}
