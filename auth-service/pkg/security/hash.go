package security

import "golang.org/x/crypto/bcrypt"

// CriptografarSenha recebe uma senha em texto plano e retorna o hash criptografado
func CriptografarSenha(senha string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// VerificarSenha compara uma senha em texto plano com um hash criptografado
func VerificarSenha(senha, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha))
	if err != nil {
		return err
	}
	return nil
}
