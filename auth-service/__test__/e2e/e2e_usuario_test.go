package e2e

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/config"
	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/router"
)

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCadastroUsuario_E2E(t *testing.T) {

	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_NAME_TEST", "")

	config.ConectarBanco(true)

	r := gin.Default()
	router.SetupRoutes(r)

	t.Run("Deve criar um usuário com sucesso", func(t *testing.T) {
		usuarioJSON := []byte(`{
            "nome":  "Cláudio Araújo",
            "email": "teste@gmail.com",
            "senha": "123456"
        }`)

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios", usuarioJSON)

		if w.Code != http.StatusCreated {
			t.Fatalf("Esperado status 201, mas obteve %d – body: %s", w.Code, w.Body.String())
		}
	
	})
}
