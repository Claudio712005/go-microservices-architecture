package security

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GerarToken gera um token JWT para o usuário com o ID fornecido
func GerarToken(usuarioID uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usuario_id": usuarioID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), 
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidarToken valida o token JWT fornecido e retorna um erro se o token for inválido
func ValidarToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func ExtrairUsuarioID(tokenString string) (uint64, error) {

	tokenString = tokenString[len("Bearer "):] 

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if usuarioID, ok := claims["usuario_id"].(float64); ok {
			return uint64(usuarioID), nil
		}
	}

	return 0, fmt.Errorf("invalid token claims")
}