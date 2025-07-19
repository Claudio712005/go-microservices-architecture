package schema

import "github.com/Claudio712005/go-microservices-architecture/auth-service/internal/domain"

type UserAuditEvent struct {
	EventType   string          `json:"event_type"`
	UserID      uint32          `json:"user_id"`
	OldUserData *domain.Usuario `json:"old_user_data,omitempty"`
	NewUserData *domain.Usuario `json:"new_user_data,omitempty"`
}
