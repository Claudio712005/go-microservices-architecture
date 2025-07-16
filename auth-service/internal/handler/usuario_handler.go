package handler

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/domain"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/mq"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/repository"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/schema"
	_ "github.com/Claudio712005/go-microservices-architecture/auth-service/internal/schema"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/error"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/response"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/security"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UsuarioHandler struct {
	repo repository.UsuarioRepository
	bus  mq.EventBus
}

func NewUsuarioHandler(repo repository.UsuarioRepository, bus mq.EventBus) *UsuarioHandler {
	return &UsuarioHandler{
		repo: repo,
		bus:  bus,
	}
}

// HandleCadastrarUsuario godoc: Cadastrar um novo usuário
// @Summary Cadastrar um novo usuário
// @Description Cadastra um novo usuário no sistema
// @Tags Usuários
// @Accept json
// @Produce json
// @Param usuario body domain.Usuario true "Dados do usuário"
// @Success 201 {object} schema.UsuarioCreatedEnvelope "Usuário cadastrado com sucesso"
// @Failure 400 {object} error.AppError "Requisição inválida"
// @Failure 409 {object} error.AppError "Conflito: usuário já existe
// @Failure 500 {object} error.AppError "Erro interno do servidor"
// @Router /usuarios [post]
// HandleCadastrarUsuario é o handler para cadastrar um usuário
func (h *UsuarioHandler) HandleCadastrarUsuario(c *gin.Context) {

	var usuario domain.Usuario
	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.Error(error.Validation("INVALID_BODY", "corpo da requisição inválido", err))
		return
	}

	if err := usuario.Validar("cadastrar"); err != nil {
		c.Error(error.Validation("INVALID_INPUT", err.Error(), err))
		return
	}

	usuarioExiste, err := h.repo.BuscarUsuarioPorEmail(usuario.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Error(error.Internal("DATABASE_ERROR", "erro ao buscar usuário por e-mail no banco de dados", err))
		return
	}

	if usuarioExiste != nil {
		c.Error(error.Conflict("USER_ALREADY_EXISTS", "já existe um usuário cadastrado com este endereço de e-mail", nil))
		return
	}

	if usuario.Senha, err = security.CriptografarSenha(usuario.Senha); err != nil {
		c.Error(error.Internal("ENCRYPTION_ERROR", "erro ao criptografar a senha do usuário", err))
		return
	}

	if usuario.ID, err = h.repo.CadastrarUsuario(usuario); err != nil {
		c.Error(error.Internal("DATABASE_ERROR", "erro ao cadastrar usuário no banco de dados", err))
		return
	}

	evt := schema.UsuarioCreated{
		ID:        usuario.ID,
		Email:     usuario.Email,
		CreatedAt: time.Now(),
	}

	if h.bus != nil {
		if err := h.bus.PublishUserCreated(c.Request.Context(), evt); err != nil {
			c.Error(error.Internal("EVENT_PUBLISH_ERROR", "usuário criado mas não foi possível publicar o evento", err))
		}
	}

	response.Created(c, fmt.Sprintf("/usuarios/%d", usuario.ID), gin.H{
		"id": usuario.ID,
	})
}

// HandleLoginUsuario godoc: Realizar login de um usuário
// @Summary Realizar login de um usuário
// @Description Realiza o login de um usuário e retorna um token JWT
// @Tags Usuários
// @Accept json
// @Produce json
// @Param login body domain.Login true "Dados de login do usuário"
// @Success 200 {object} schema.LoginResponseEnvelope "Login realizado com sucesso"
// @Failure 400 {object} error.AppError "Requisição inválida"
// @Failure 401 {object} error.AppError "Credenciais inválidas"
// @Failure 404 {object} error.AppError "Usuário não encontrado"
// @Failure 500 {object} error.AppError "Erro interno do servidor"
// @Router /usuarios/login [post]
// HandleLoginUsuario é o handler para realizar o login de um usuário
func (h *UsuarioHandler) HandleLoginUsuario(c *gin.Context) {
	var login domain.Login

	if err := c.ShouldBindJSON(&login); err != nil {
		c.Error(error.Validation("INVALID_BODY", "corpo da requisição inválido", err))
		return
	}

	if err := login.Validar(); err != nil {
		c.Error(error.Validation("INVALID_INPUT", err.Error(), err))
		return
	}

	usuario, err := h.repo.BuscarUsuarioPorEmail(login.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(error.NotFound("USER_NOT_FOUND", "usuário não encontrado com este e-mail", nil))
			return
		}
		c.Error(error.Internal("DATABASE_ERROR", "erro ao buscar usuário no banco de dados", err))
		return
	}

	if err := security.VerificarSenha(login.Senha, usuario.Senha); err != nil {
		c.Error(error.Unauthorized("INVALID_CREDENTIALS", "credenciais inválidas", nil))
		return
	}

	token, err := security.GerarToken(usuario.ID)
	if err != nil {
		c.Error(error.Internal("TOKEN_GENERATION_ERROR", "erro ao gerar token JWT", err))
		return
	}

	usuario.Senha = ""

	response.OK(c, gin.H{
		"token":   token,
		"usuario": usuario,
	})
}

// HandleBuscarUsuarioLogado godoc: Buscar informações do usuário logado
// @Summary Buscar informações do usuário logado
// @Description Busca as informações do usuário logado a partir do token JWT
// @Tags Usuários
// @Accept json
// @Produce json
// @Success 200 {object} schema.UsuarioEnvelope "Usuário encontrado"
// @Failure 401 {object} error.AppError "Token inválido ou expirado"
// @Failure 404 {object} error.AppError "Usuário não encontrado"
// @Failure 500 {object} error.AppError "Erro interno do servidor"
// @Router /usuarios/logado [get]
// HandleBuscarUsuarioLogado é o handler para buscar informações do usuário logado
func (h *UsuarioHandler) HandleBuscarUsuarioLogado(c *gin.Context) {
	idToken, err := security.ExtrairUsuarioID(c.GetHeader("Authorization"))
	if err != nil {
		c.Error(error.Unauthorized("INVALID_TOKEN", "token inválido ou expirado", err))
		return
	}

	usuario, err := h.repo.BuscarUsuarioPorID(idToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(error.NotFound("USER_NOT_FOUND", "usuário não encontrado com este ID", nil))
			return
		}
		c.Error(error.Internal("DATABASE_ERROR", "erro ao buscar usuário no banco de dados", err))
		return
	}

	response.OK(c, usuario)
}

// HandleAlterarUsuario godoc: Alterar informações do usuário
// @Summary Alterar informações do usuário
// @Description Altera as informações do usuário logado
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Param usuario body domain.Usuario true "Dados do usuário"
// @Success 200 {object} schema.UsuarioEnvelope "Usuário atualizado com sucesso"
// @Failure 400 {object} error.AppError "Requisição inválida"
// @Failure 401 {object} error.AppError "Token inválido ou expirado"
// @Failure 403 {object} error.AppError "Acesso negado"
// @Failure 404 {object} error.AppError "Usuário não encontrado"
// @Failure 500 {object} error.AppError "Erro interno do servidor"
// @Router /usuarios/{id} [put]
// HandleAlterarUsuario é o handler para alterar as informações de um usuário
func (h *UsuarioHandler) HandleAlterarUsuario(c *gin.Context) {
	id := c.Param("id")

	idUsuario, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.Error(error.Validation("INVALID_ID", "ID do usuário inválido", err))
		return
	}

	tokenID, err := security.ExtrairUsuarioID(c.GetHeader("Authorization"))
	if err != nil {
		c.Error(error.Unauthorized("INVALID_TOKEN", "token inválido ou expirado", err))
		return
	}

	if tokenID != uint32(idUsuario) {
		c.Error(error.Forbidden("FORBIDDEN", "você não tem permissão para alterar este usuário", nil))
		return
	}

	var usuario domain.Usuario
	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.Error(error.Validation("INVALID_BODY", "corpo da requisição inválido", err))
		return
	}

	usuario.ID = uint32(idUsuario)

	if err := usuario.Validar("atualizar"); err != nil {
		c.Error(error.Validation("INVALID_INPUT", err.Error(), err))
		return
	}

	usuarioExistente, err := h.repo.BuscarUsuarioPorEmail(usuario.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		c.Error(error.Internal("DATABASE_ERROR", "erro ao buscar usuário por e-mail no banco de dados", err))
		return
	}

	if usuarioExistente != nil && usuarioExistente.ID != usuario.ID {
		c.Error(error.Conflict("USER_ALREADY_EXISTS", "já existe um usuário cadastrado com este endereço de e-mail", nil))
		return
	}

	usuarioSalvo, err := h.repo.EditarUsuario(usuario)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(error.NotFound("USER_NOT_FOUND", "usuário não encontrado com este ID", nil))
			return
		}
		c.Error(error.Internal("DATABASE_ERROR", "erro ao atualizar usuário no banco de dados", err))
		return
	}

	usuarioSalvo.Senha = ""

	response.OK(c, *usuarioSalvo)
}

// HandleAlterarSenhaUsuarioLogado godoc: Alterar senha do usuário logado
// @Summary Alterar senha do usuário logado
// @Description Altera a senha do usuário logado
// @Tags Usuários
// @Accept json
// @Produce json
// @Param senha body domain.Senha true "Nova senha do usuário"
// @Success 200 {object} schema.MessageEnvelope "Senha atualizada com sucesso"
// @Failure 400 {object} error.AppError "Requisição inválida"
// @Failure 401 {object} error.AppError "Token inválido ou expirado"
// @Failure 403 {object} error.AppError "Acesso negado"
// @Failure 404 {object} error.AppError "Usuário não encontrado"
// @Failure 500 {object} error.AppError "Erro interno do servidor"
// @Router /usuarios/logado/senha [put]
// HandleAlterarSenhaUsuarioLogado é o handler para alterar a senha do usuário logado
func (h *UsuarioHandler) HandleAlterarSenhaUsuarioLogado(c *gin.Context) {
	idToken, err := security.ExtrairUsuarioID(c.GetHeader("Authorization"))
	if err != nil {
		c.Error(error.Unauthorized("INVALID_TOKEN", "token inválido ou expirado", err))
	}

	senhaHash, err := h.repo.BuscarSenha(idToken)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(error.NotFound("USER_NOT_FOUND", "usuário não encontrado com este ID", nil))
			return
		}
		c.Error(error.Internal("DATABASE_ERROR", "erro ao buscar usuário no banco de dados", err))
		return
	}

	var senha domain.Senha
	if err := c.ShouldBindJSON(&senha); err != nil {
		c.Error(error.Validation("INVALID_BODY", "corpo da requisição inválido", err))
		return
	}

	if err := senha.Validar(); err != nil {
		c.Error(error.Validation("INVALID_INPUT", err.Error(), err))
		return
	}

	if err := senha.ValidarSenha(senhaHash); err != nil {
		c.Error(error.Forbidden("INVALID_PASSWORD", "senha incorretas", err))
		return
	}

	senha.SenhaNova, err = security.CriptografarSenha(senha.SenhaNova)
	if err != nil {
		c.Error(error.Internal("ENCRYPTION_ERROR", "erro ao criptografar a nova senha", err))
		return
	}

	if _, err := h.repo.AtualizarSenha(domain.Usuario{
		ID:    idToken,
		Senha: senha.SenhaNova,
	}); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(error.NotFound("USER_NOT_FOUND", "usuário não encontrado com este ID", nil))
			return
		}
		c.Error(error.Internal("DATABASE_ERROR", "erro ao atualizar senha no banco de dados", err))
		return
	}

	response.OK(c, gin.H{
		"message": "Senha atualizada com sucesso",
	})

}

// HandleDeletarUsuario godoc: Deletar um usuário
// @Summary Deletar um usuário
// @Description Deleta um usuário do sistema
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path string true "ID do usuário"
// @Success 204 {object} interface{} "Usuário deletado com sucesso"
// @Failure 400 {object} error.AppError "Requisição inválida"
// @Failure 401 {object} error.AppError "Token inválido ou expirado"
// @Failure 403 {object} error.AppError "Acesso negado"
// @Failure 404 {object} error.AppError "Usuário não encontrado"
// @Failure 500 {object} error.AppError "Erro interno do servidor"
// @Router /usuarios/{id} [delete]
// HandleDeletarUsuario é o handler para deletar um usuário
func (h *UsuarioHandler) HandleDeletarUsuario(c *gin.Context) {

	id := c.Param("id")

	idUsuario, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.Error(error.Validation("INVALID_ID", "ID do usuário inválido", err))
		return
	}

	tokenID, err := security.ExtrairUsuarioID(c.GetHeader("Authorization"))
	if err != nil {
		c.Error(error.Unauthorized("INVALID_TOKEN", "token inválido ou expirado", err))
		return
	}

	if tokenID != uint32(idUsuario) {
		c.Error(error.Forbidden("FORBIDDEN", "você não tem permissão para deletar este usuário", nil))
		return
	}

	if err := h.repo.DeletarUsuario(uint32(idUsuario)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Error(error.NotFound("USER_NOT_FOUND", "usuário não encontrado com este ID", nil))
			return
		}
		c.Error(error.Internal("DATABASE_ERROR", "erro ao deletar usuário no banco de dados", err))
		return
	}

	response.NoContent(c)
}
