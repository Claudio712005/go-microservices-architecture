package domain

import "github.com/Claudio712005/go-microservices-architecture/auth-service/pkg"

type Login struct {
	Email string `json:"email" binding:"required,email"`
	Senha string `json:"senha" binding:"required,min=6"`
}

func (l *Login) Validar() error {
	if err := pkg.ValidarCampos(l); err != nil {
		return err
	}

	return nil
}