package repository

import (
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/domain"
	"gorm.io/gorm"
)

// UsuarioRepository define os métodos para interagir com o repositório de usuários
type UsuarioRepository interface {
	CadastrarUsuario(usuario domain.Usuario) (uint32, error)
	BuscarUsuarioPorEmail(email string) (*domain.Usuario, error)
	BuscarUsuarioPorID(id uint32) (*domain.Usuario, error)
	EditarUsuario(usuario domain.Usuario) (*domain.Usuario, error)
	AtualizarSenha(usuario domain.Usuario) (*domain.Usuario, error)
	BuscarSenha(id uint32) (string, error)
}

type usuarioRepository struct {
	db *gorm.DB
}

// NewUsuarioRepository cria uma nova instância de UsuarioRepository
func NewUsuarioRepository(db *gorm.DB) UsuarioRepository {
	return &usuarioRepository{db: db}
}

// CadastrarUsuario cadastra um novo usuário no repositório
func (r *usuarioRepository) CadastrarUsuario(usuario domain.Usuario) (uint32, error) {
	if err := r.db.Create(&usuario).Error; err != nil {
		return 0, err
	}
	return usuario.ID, nil
}

// BuscarUsuarioPorEmail busca um usuário pelo email no repositório
func (r *usuarioRepository) BuscarUsuarioPorEmail(email string) (*domain.Usuario, error) {
	var usuario domain.Usuario

	if err := r.db.Where("email = ?", email).First(&usuario).Error; err != nil {
		return nil, err
	}

	return &usuario, nil
}

// BuscarUsuarioPorID busca um usuário pelo ID no repositório
func (r *usuarioRepository) BuscarUsuarioPorID(id uint32) (*domain.Usuario, error) {
	var usuario domain.Usuario

	if err := r.db.Select("id, nome, email").First(&usuario, id).Error; err != nil {
		return nil, err
	}

	return &usuario, nil
}

// EditarUsuario atualiza os dados de um usuário no repositório
func (r *usuarioRepository) EditarUsuario(usuario domain.Usuario) (*domain.Usuario, error) {
	if usuario.ID == 0 {
		return &domain.Usuario{}, gorm.ErrRecordNotFound
	}

	update := map[string]interface{}{
		"nome":  usuario.Nome,
		"email": usuario.Email,
	}

	if err := r.db.Model(&domain.Usuario{}).Where("id = ?", usuario.ID).Updates(update).Error; err != nil {
		return &domain.Usuario{}, err
	}


	return &usuario, nil
}

// AtualizarSenha atualiza a senha de um usuário no repositório
func (r *usuarioRepository) AtualizarSenha(usuario domain.Usuario) (*domain.Usuario, error) {
	if usuario.ID == 0 {
		return &domain.Usuario{}, gorm.ErrRecordNotFound
	}

	if err := r.db.Model(&domain.Usuario{}).Where("id = ?", usuario.ID).Update("senha", usuario.Senha).Error; err != nil {
		return &domain.Usuario{}, err
	}

	return &usuario, nil
}

// BuscarSenha busca a senha de um usuário pelo ID
func (r *usuarioRepository) BuscarSenha(id uint32) (string, error) {
	var usuario domain.Usuario
	if err := r.db.Select("senha").First(&usuario, id).Error; err != nil {
		return "", err
	}

	return usuario.Senha, nil
}