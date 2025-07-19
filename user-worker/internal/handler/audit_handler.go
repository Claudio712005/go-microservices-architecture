package handler

import (
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/error"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/response"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/security"
	"github.com/Claudio712005/go-microservices-architecture/user-worker/internal/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditHandler struct {
	repo repository.AuditRepository
}

func NewAuditHandler(repo repository.AuditRepository) *AuditHandler {
	return &AuditHandler{
		repo: repo,
	}
}

// HandleLogsUsuarioLogado é o manipulador para recuperar os logs de auditoria do usuário logado.
func (h *AuditHandler) HandleLogsUsuarioLogado(c *gin.Context) {

	idToken, err := security.ExtrairUsuarioID(c.GetHeader("Authorization"))
	if err != nil {
		c.Error(error.Unauthorized("INVALID_TOKEN", "token inválido ou expirado", err))
		return
	}

	audits, err := h.repo.GetAuditByUserID(idToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(error.NotFound("AUDIT_NOT_FOUND", "nenhum evento de auditoria encontrado para o usuário", err))
			return
		}
		c.Error(error.Internal("INTERNAL_ERROR", "erro ao recuperar eventos de auditoria", err))
		return
	}

	response.OK(c, audits)
}