package domain

import (
	"time"
)

// AuditEventsChanges representa as mudanças em um evento de auditoria
type AuditEventsChanges struct {
	ID           uint32      `gorm:"primaryKey" json:"id"`
	AuditEventID uint32      `gorm:"not null;index" json:"audit_event_id"`
	FieldName    string      `gorm:"size:100;not null" json:"field_name"`
	OldValue     string      `gorm:"type:text" json:"old_value"`
	NewValue     string      `gorm:"type:text" json:"new_value"`

	AuditEvent   AuditEvent  `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}

// AuditEvent representa um evento de auditoria
type AuditEvent struct {
	ID        uint32               `gorm:"primaryKey" json:"id"`
	UserID    uint32               `gorm:"not null;index" json:"user_id"`
	EventType string               `gorm:"size:100;not null" json:"event_type"`
	Timestamp time.Time            `gorm:"not null" json:"timestamp"`
	Source    string               `gorm:"size:255;not null" json:"source"`

	Changes   []AuditEventsChanges `gorm:"foreignKey:AuditEventID;constraint:OnDelete:CASCADE;" json:"changes"`
	User      Usuario              `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}

// Usuario representa a tabela de usuários
type Usuario struct {
	ID     uint32        `gorm:"primaryKey" json:"id"`
	Nome   string        `gorm:"size:100;not null" json:"nome"`
	Email  string        `gorm:"size:255;unique;not null" json:"email"`
	Senha  string        `gorm:"size:255;not null" json:"senha"`

	AuditEvents []AuditEvent `gorm:"foreignKey:UserID" json:"-"`
}
