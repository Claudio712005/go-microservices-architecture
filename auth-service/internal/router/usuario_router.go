package router

import (
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/config"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/handler"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/mq"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/repository"
)

func getUsuarioRoutes(
	bus mq.EventBus,
) []Route {

	h := handler.NewUsuarioHandler(repository.NewUsuarioRepository(config.DB), bus)	

	return []Route{
		{
			Path:        "",
			Method:      "POST",
			HandlerFunc: h.HandleCadastrarUsuario,
			HasAuth:     false,
		},
		{
			Path: "/login",
			Method:      "POST",
			HandlerFunc: h.HandleLoginUsuario,
			HasAuth:     false,
		},
		{
			Path: "/logado",
			Method:      "GET",
			HandlerFunc: h.HandleBuscarUsuarioLogado,
			HasAuth:     true,
		},
		{
			Path: "/:id",
			Method:      "PUT",
			HandlerFunc: h.HandleAlterarUsuario,
			HasAuth:     true,
		}, 
		{
			Path: "/senha",
			Method:      "PUT",
			HandlerFunc: h.HandleAlterarSenhaUsuarioLogado,
			HasAuth:     true,
		},
		{
			Path: "/:id",
			Method:      "DELETE",
			HandlerFunc: h.HandleDeletarUsuario,
			HasAuth:     true,
		},
	}
}
