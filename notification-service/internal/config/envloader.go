package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv carrega as variáveis de ambiente a partir de arquivos .env
func LoadEnv() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Erro ao carregar arquivo .env: %v", err)
	}
}
