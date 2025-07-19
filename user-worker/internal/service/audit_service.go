package service

import (
	"errors"
	"log"
	"time"

	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/domain"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/repository"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/schema"
)

// AuditService providencia serviços relacionados a auditoria.
type AuditService struct {
	repo repository.AuditRepository
}

// NewAuditService cria uma nova instância de AuditService.
func NewAuditService(repo repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (a *AuditService) RegistryNewAudit(event *schema.UserAuditEvent) error {
	
	if event == nil {
		log.Println("RegistryNewAudit: event is nil")
		return errors.New("event cannot be nil")
	}

	log.Println("RegistryNewAudit: registering new audit event body: ", event)

	var auditEvent domain.AuditEvent

	auditEvent.UserID = event.UserID
	auditEvent.EventType = event.EventType
	auditEvent.Timestamp = time.Now()
	auditEvent.User = domain.Usuario{
		ID: event.UserID,
	}

	auditEvent.Source = "user-worker"

	if event.OldUserData == nil {
		auditEvent.Changes = []domain.AuditEventsChanges{
			{
				FieldName: "nome",
				OldValue:  "",
				NewValue:  event.NewUserData.Nome,
			},
			{
				FieldName: "email",
				OldValue:  "",
				NewValue:  event.NewUserData.Email,
			},
			{
				FieldName: "senha",
				OldValue:  "",
				NewValue:  event.NewUserData.Senha,
			},
		}
	} else {
		if event.OldUserData.Nome != event.NewUserData.Nome {
			auditEvent.Changes = append(auditEvent.Changes, domain.AuditEventsChanges{
				FieldName: "nome",
				OldValue:  event.OldUserData.Nome,
				NewValue:  event.NewUserData.Nome,
			})
		}
		if event.OldUserData.Email != event.NewUserData.Email {
			auditEvent.Changes = append(auditEvent.Changes, domain.AuditEventsChanges{
				FieldName: "email",
				OldValue:  event.OldUserData.Email,
				NewValue:  event.NewUserData.Email,
			})
		}
		if event.OldUserData.Senha != event.NewUserData.Senha {
			auditEvent.Changes = append(auditEvent.Changes, domain.AuditEventsChanges{
				FieldName: "senha",
				OldValue:  event.OldUserData.Senha,
				NewValue:  event.NewUserData.Senha,
			})
		}
	}

	if err := a.repo.RegistryNewAudit(&auditEvent); err != nil {
		log.Printf("RegistryNewAudit: error saving audit event: %v", err)
		return err
	}

	return nil
}