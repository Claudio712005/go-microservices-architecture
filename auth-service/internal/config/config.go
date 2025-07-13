// internal/config/config.go
package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Load() {
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

	if env == "test" {
		if testDB := os.Getenv("DB_NAME_TEST"); testDB != "" {
			_ = os.Setenv("DB_NAME", testDB)
		}
	}

	ConectarBanco()
}
