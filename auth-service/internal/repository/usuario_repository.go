package repository

import "gorm.io/gorm"

type UsuarioRepository interface {
}

type usuarioRepository struct {
	db *gorm.DB
}

func NewUsuarioRepository(db *gorm.DB) UsuarioRepository {
	return &usuarioRepository{db: db}
}
