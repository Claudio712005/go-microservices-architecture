package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func findRepoRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// LoadEnv carrega as variáveis de ambiente a partir de arquivos .env
func LoadEnv() {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	root := findRepoRoot()
	paths := []string{
		filepath.Join(root, ".env."+env),
		filepath.Join(root, ".env"),
	}

	_ = godotenv.Load(paths...)
}
