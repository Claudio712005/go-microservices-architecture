package repository

import (
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/domain"
	"gorm.io/gorm"
)

// AuditRepository é a estrutura que implementa o repositório de auditoria.
type AuditRepository interface {
	RegistryNewAudit(event *domain.AuditEvent) error
}

type auditRepository struct {
	db *gorm.DB
}

// NewAuditRepository cria uma nova instância de AuditRepository.
func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

// RegistryNewAudit registra um novo evento de auditoria no banco de dados.
func (r *auditRepository) RegistryNewAudit(event *domain.AuditEvent) error {
	if err := r.db.Create(event).Error; err != nil {
		return err
	}
	return nil
}