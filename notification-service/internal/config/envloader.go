package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

// LoadEnv carrega as variáveis de ambiente a partir de arquivos .env
func LoadEnv() {

	if err := godotenv.Load(); err != nil {
		fmt.Printf("Erro ao carregar arquivo .env: %v", err)
	}
}
