package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/domain"
)

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCadastroUsuario_E2E(t *testing.T) {
	t.Run("Cadastrar usuário com dados válidos", func(t *testing.T) {
		usuarioJSON := []byte(`{
            "nome":  "Cláudio Araújo",
            "email": "teste@gmail.com",
            "senha": "123456"
        }`)

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios", usuarioJSON)

		if w.Code != http.StatusCreated {
			t.Fatalf("Esperado status 201, obtido %d – body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Data struct {
				ID uint32 `json:"id"`
			} `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Falha ao decodificar resposta JSON: %v", err)
		}

		if resp.Data.ID == 0 {
			t.Fatalf("ID do usuário retornado é inválido")
		}
	})

	t.Run("Cadastrar usuário com email já existente", func(t *testing.T) {
		usuarioJSON := []byte(`{
			"nome":  "Cláudio Araújo",
			"email": "teste@gmail.com",
			"senha": "123456"
		}`)
		w := performRequest(r, http.MethodPost, "/api/v1/usuarios", usuarioJSON)

		if w.Code != http.StatusConflict {
			t.Fatalf("Esperado status 409, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Cadastrar usuário com dados inválidos", func(t *testing.T) {
		usuarioJSON := []byte(`{
			"nome":  "",
			"email": "email-invalido",
			"senha": "123"
		}`)

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios", usuarioJSON)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperado status 400, obtido %d – body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Message string `json:"message"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Falha ao decodificar resposta JSON: %v", err)
		}

		if resp.Message == "" {
			t.Fatal("Mensagem de erro não foi retornada")
		}
	})
}

func TestLoginUsuario_E2E(t *testing.T) {

	t.Run("Login com dados válidos", func(t *testing.T) {
		loginJSON := []byte(`{
			"email": "teste@gmail.com",
			"senha": "123456"
		}`)

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios/login", loginJSON)

		if w.Code != http.StatusOK {
			t.Fatalf("Esperado status 200, obtido %d – body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Data struct {
				Token   string         `json:"token"`
				Usuario domain.Usuario `json:"usuario"`
			} `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Falha ao decodificar resposta JSON: %v", err)
		}

		if resp.Data.Token == "" {
			t.Fatal("Token não foi retornado")
		}

		if resp.Data.Usuario.ID == 0 {
			t.Fatal("ID do usuário não foi retornado")
		}

		if resp.Data.Usuario.Senha != "" {
			t.Fatal("Senha do usuário não deve ser retornada no login")
		}
	})

	t.Run("Login com email não cadastrado", func(t *testing.T) {
		loginJSON := []byte(`{
			"email": "emailInvalido@gmail.com",
			"senha": "123456"
		}`)

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios/login", loginJSON)
		if w.Code != http.StatusNotFound {
			t.Fatalf("Esperado status 404, obtido %d – body: %s", w.Code, w.Body.String())
		}

	})

	t.Run("Login com senha incorreta", func(t *testing.T) {
		loginJSON := []byte(`{
			"email": "teste@gmail.com",
			"senha": "senhaIncorreta"
		}`)
		w := performRequest(r, http.MethodPost, "/api/v1/usuarios/login", loginJSON)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperado status 401, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Login com dados inválidos", func(t *testing.T) {
		loginJSON := []byte(`{
			"email": "emailInvalido",
			"senha": "123"
		}`)

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios/login", loginJSON)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Esperado status 400, obtido %d – body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Message string `json:"message"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Falha ao decodificar resposta JSON: %v", err)
		}

		if resp.Message == "" {
			t.Fatal("Mensagem de erro não foi retornada")
		}
	})
}
