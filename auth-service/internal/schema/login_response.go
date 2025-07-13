package schema

import "github.com/Claudio712005/go-microservices-architecture/auth-service/internal/domain"

// LoginResponseEnvelope godoc: Resposta do login do usuário
// LoginResponseEnvelope é a estrutura de resposta para o login do usuário
// @Description Resposta do login do usuário
// @name LoginResponseEnvelope
type LoginResponseEnvelope struct {
	Data struct {
		Token   string `json:"token"`
		Usuario domain.Usuario `json:"usuario"`
	}
}