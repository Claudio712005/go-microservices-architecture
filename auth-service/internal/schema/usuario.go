package schema

import "github.com/Claudio712005/go-microservices-architecture/auth-service/internal/domain"

// UsuarioCreatedEnvelope godoc: Resposta do usuário criado
// UsuarioCreatedEnvelope é a estrutura de resposta para o usuário criado
// @Description Resposta do usuário criado
// @name UsuarioCreatedEnvelope
type UsuarioEnvelope struct {
	Data struct {
		Usuario domain.Usuario `json:"usuario"`
	}
}