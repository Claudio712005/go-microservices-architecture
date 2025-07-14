package domain

import (
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/security"
)

type Senha struct {
	SenhaNova  string `json:"senha_nova" validate:"required,min=6,max=100"`
	SenhaAtual string `json:"senha_atual" validate:"required,min=6,max=100"`
}

func (s *Senha) Validar() error {

	if err := pkg.ValidarCampos(s); err != nil {
		return err
	}

	return nil
}

func (s *Senha) ValidarSenha(senhaAtualHash string) error {

	if err := security.VerificarSenha(s.SenhaAtual, senhaAtualHash) ; err != nil {
		return err
	}

	return nil
}