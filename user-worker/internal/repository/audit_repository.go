package repository

import (
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/domain"
	"gorm.io/gorm"
)

// AuditRepository é a estrutura que implementa o repositório de auditoria.
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository cria uma nova instância de AuditRepository.
func NewAudditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// RegistryNewAudit registra um novo evento de auditoria no banco de dados.
func (r *AuditRepository) RegistryNewAudit(event *domain.AuditEvent) error {
	if err := r.db.Create(event).Error; err != nil {
		return err
	}
	return nil
}