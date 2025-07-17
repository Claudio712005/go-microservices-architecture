package domain

import "time"

// Change representa uma alteração em um evento de auditoria
type Change struct {
	Field     string `json:"field"`
	OldValue  string `json:"old_value"`
	NewValue string `json:"new_value"`
}

// AuditEvent representa um evento de auditoria
type AuditEvent struct {
	UserID    uint32 `json:"user_id"`
	EventType string `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Source string `json:"source"`
	Changes []Change `json:"changes"`
}