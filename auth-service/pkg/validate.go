package pkg

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidarCampos valida os campos de uma estrutura e retorna mensagens de erro amigáveis
func ValidarCampos(i interface{}) error {
	validate := validator.New()
	if err := validate.Struct(i); err != nil {
		var mensagens []string
		for _, fieldErr := range err.(validator.ValidationErrors) {
			mensagens = append(mensagens, traduzirErro(fieldErr))
		}
		return fmt.Errorf("Campos incorretos: %s", strings.Join(mensagens, " | "))
	}
	return nil
}

func traduzirErro(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("O campo '%s' é obrigatório.", fe.Field())
	case "min":
		return fmt.Sprintf("O campo '%s' deve ter no mínimo %s caracteres.", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("O campo '%s' deve ter no máximo %s caracteres.", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("O campo '%s' deve ser um e-mail válido.", fe.Field())
	default:
		return fmt.Sprintf("O campo '%s' é inválido (regra: %s).", fe.Field(), fe.Tag())
	}
}
