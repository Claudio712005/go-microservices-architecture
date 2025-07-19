package router

import (
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/config"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/handler"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/repository"
)

func getAuditedUsuarioRoutes() []Route {
	h := handler.NewAuditHandler(repository.NewAuditRepository(config.DB))

	return []Route{
		{
			Path: 	  "/logado",
			Method:   "GET",
			HandlerFunc: h.HandleLogsUsuarioLogado,
			HasAuth:  true,
		},
	}
}