package schema

import "time"

// UsuarioCreatedEnvelope godoc: Resposta do usuário criado
// UsuarioCreatedEnvelope é a estrutura de resposta para o usuário criado
// @Description Resposta do usuário criado
// @name UsuarioCreatedEnvelope
type UsuarioCreatedEnvelope struct {
	Data    struct {
		ID uint32 `json:"id"`
	} `json:"data"`
}

// UsuarioCreated é a estrutura que representa o usuário criado
type UsuarioCreated struct {
	ID uint32 `json:"id"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}