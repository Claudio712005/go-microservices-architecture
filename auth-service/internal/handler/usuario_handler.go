package handler

import (
	"fmt"
    _ "github.com/Claudio712005/go-microservices-architecture/auth-service/internal/schema"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/config"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/domain"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/repository"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/error"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/response"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/pkg/security"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
func HandleCadastrarUsuario(c *gin.Context) {

	var usuario domain.Usuario
	if err := c.ShouldBindJSON(&usuario); err != nil {
		c.Error(error.Validation("INVALID_BODY", "corpo da requisição inválido", err))
		return
	}

	if err := usuario.Validar("cadastrar"); err != nil {
		c.Error(error.Validation("INVALID_INPUT", err.Error(), err))
		return
	}

	repositorio := repository.NewUsuarioRepository(config.DB)

	usuarioExiste, err := repositorio.BuscarUsuarioPorEmail(usuario.Email)
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

	if usuario.ID, err = repositorio.CadastrarUsuario(usuario); err != nil {
		c.Error(error.Internal("DATABASE_ERROR", "erro ao cadastrar usuário no banco de dados", err))
		return
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
func HandleLoginUsuario(c *gin.Context) {
	var login domain.Login

	if err := c.ShouldBindJSON(&login); err != nil {
		c.Error(error.Validation("INVALID_BODY", "corpo da requisição inválido", err))
		return
	}

	if err := login.Validar(); err != nil {
		c.Error(error.Validation("INVALID_INPUT", err.Error(), err))
		return
	}

	repositorio := repository.NewUsuarioRepository(config.DB)
	usuario, err := repositorio.BuscarUsuarioPorEmail(login.Email)
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
		"token": token,
		"usuario": usuario,
	})
}