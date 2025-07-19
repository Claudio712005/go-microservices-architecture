package service

import (
	"testing"

	"github.com/Claudio712005/go-microservices-architecture/user-worker/__tests__/mocks"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/domain"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/schema"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/service"
	"go.uber.org/mock/gomock"
)

func TestAuditService_RegistryAudit_EventNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockAuditRepository(ctrl)
	auditService := service.NewAuditService(repo)

	err := auditService.RegistryNewAudit(nil)

	if err == nil || err.Error() != "event cannot be nil" {
		t.Errorf("esperado erro 'event cannot be nil', mas recebeu: %v", err)
	}
}

func TestAuditService_RegistryAudit_NewUserDataOnly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockAuditRepository(ctrl)
	auditService := service.NewAuditService(repo)

	event := schema.UserAuditEvent{
		UserID:    1,
		EventType: "user.created",
		NewUserData: &domain.Usuario{
			Nome:  "Novo",
			Email: "novo@teste.com",
			Senha: "123",
		},
		OldUserData: nil,
	}

	repo.EXPECT().RegistryNewAudit(gomock.Any()).DoAndReturn(func(audit *domain.AuditEvent) error {
		if len(audit.Changes) != 3 {
			t.Errorf("esperado 3 mudanças, recebeu %d", len(audit.Changes))
		}
		return nil
	}).Times(1)

	err := auditService.RegistryNewAudit(&event)

	if err != nil {
		t.Errorf("esperado sucesso, mas recebeu erro: %v", err)
	}
}

func TestAuditService_RegistryAudit_ChangedOnlyNome(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockAuditRepository(ctrl)
	auditService := service.NewAuditService(repo)

	event := schema.UserAuditEvent{
		UserID:    1,
		EventType: "user.updated",
		NewUserData: &domain.Usuario{
			Nome:  "Novo Nome",
			Email: "igual@teste.com",
			Senha: "senha123",
		},
		OldUserData: &domain.Usuario{
			Nome:  "Antigo Nome",
			Email: "igual@teste.com",
			Senha: "senha123",
		},
	}

	repo.EXPECT().RegistryNewAudit(gomock.Any()).DoAndReturn(func(audit *domain.AuditEvent) error {
		if len(audit.Changes) != 1 {
			t.Errorf("esperado 1 mudança, recebeu %d", len(audit.Changes))
		}
		if audit.Changes[0].FieldName != "nome" {
			t.Errorf("esperado mudança no campo 'nome', recebeu '%s'", audit.Changes[0].FieldName)
		}
		return nil
	}).Times(1)

	err := auditService.RegistryNewAudit(&event)

	if err != nil {
		t.Errorf("esperado sucesso, mas recebeu erro: %v", err)
	}
}
