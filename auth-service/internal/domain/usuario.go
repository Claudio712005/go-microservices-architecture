package domain

import "github.com/Claudio712005/go-microservices-architecture/auth-service/pkg"

// Usuario representa a estrutura de um usuário no sistema
type Usuario struct {
	ID    uint32 `json:"id" gorm:"primaryKey"`
	Nome  string `json:"nome" gorm:"not null" validate:"required,min=3,max=100"`
	Email string `json:"email" gorm:"not null;unique" validate:"required,email"`
	Senha string `json:"senha" gorm:"not null" validate:"required,min=6,max=100"`
}

// Validar valida os campos do usuário
func (u *Usuario) Validar(tipo string) error {

	if err := pkg.ValidarCampos(u); err != nil && tipo == "cadastrar" {
		return err
	}

	return nil
}
