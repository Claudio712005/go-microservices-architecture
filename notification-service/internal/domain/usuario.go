package domain

// Usuario representa um usuário no sistema.
type Usuario struct {
	ID        uint32 `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}
