package config

import (
	"os"
	"path/filepath"
)

// findRepoRoot sobe diretórios até achar um arquivo‑marcador (go.mod, .git).
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
