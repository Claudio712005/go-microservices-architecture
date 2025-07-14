package schema

// UsuarioEnvelope é a estrutura de resposta para o usuário criado
// UsuarioEnvelope godoc: Resposta do usuário criado
// @Description Resposta do usuário criado
// @name UsuarioCreatedEnvelope
type MessageEnvelope struct {
	Data struct {
		Message string `json:"message"`
	} `json:"data"`
}