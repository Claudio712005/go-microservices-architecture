package schema

// UsuarioCreatedEnvelope godoc: Resposta do usuário criado
// UsuarioCreatedEnvelope é a estrutura de resposta para o usuário criado
// @Description Resposta do usuário criado
// @name UsuarioCreatedEnvelope
type UsuarioCreatedEnvelope struct {
	Data    struct {
		ID uint32 `json:"id"`
	} `json:"data"`
}
