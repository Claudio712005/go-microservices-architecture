package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Claudio712005/go-microservices-architecture/auth-service/internal/domain"
)

func performRequest(r http.Handler, method, path string, body []byte, header http.Header) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

var token string

func TestCadastroUsuario_E2E(t *testing.T) {
	t.Run("Cadastrar usuário com dados válidos", func(t *testing.T) {
		usuarioJSON := []byte(`{
            "nome":  "Cláudio Araújo",
            "email": "teste@gmail.com",
            "senha": "123456"
        }`)

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios", usuarioJSON, nil)

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
		w := performRequest(r, http.MethodPost, "/api/v1/usuarios", usuarioJSON, nil)

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

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios", usuarioJSON, nil)

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

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios/login", loginJSON, nil)

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

		token = resp.Data.Token

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

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios/login", loginJSON, nil)
		if w.Code != http.StatusNotFound {
			t.Fatalf("Esperado status 404, obtido %d – body: %s", w.Code, w.Body.String())
		}

	})

	t.Run("Login com senha incorreta", func(t *testing.T) {
		loginJSON := []byte(`{
			"email": "teste@gmail.com",
			"senha": "senhaIncorreta"
		}`)
		w := performRequest(r, http.MethodPost, "/api/v1/usuarios/login", loginJSON, nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperado status 401, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Login com dados inválidos", func(t *testing.T) {
		loginJSON := []byte(`{
			"email": "emailInvalido",
			"senha": "123"
		}`)

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios/login", loginJSON, nil)

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

func TestBuscarUsuarioLogado(t *testing.T) {

	t.Run("Buscar usuário logado com token válido", func(t *testing.T) {
		w := performRequest(r, http.MethodGet, "/api/v1/usuarios/logado", nil, http.Header{
			"Authorization": []string{"Bearer " + token},
		})

		if w.Code != http.StatusOK {
			t.Fatalf("Esperado status 200, obtido %d – body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Data domain.Usuario `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Falha ao decodificar resposta JSON: %v", err)
		}

		if resp.Data.ID == 0 {
			t.Fatal("ID do usuário não foi retornado")
		}
		if resp.Data.Senha != "" {
			t.Fatal("Senha do usuário não deve ser retornada")
		}
	})

	t.Run("Buscar usuário logado com token inválido", func(t *testing.T) {
		w := performRequest(r, http.MethodGet, "/api/v1/usuarios/logado", nil, http.Header{
			"Authorization": []string{"Bearer tokenInvalido"},
		})
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperado status 401, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Buscar usuário logado sem token", func(t *testing.T) {
		w := performRequest(r, http.MethodGet, "/api/v1/usuarios/logado", nil, nil)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperado status 401, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})
}

func TestAlterarUsuario(t *testing.T) {
	t.Run("Alterar usuário com dados válidos", func(t *testing.T) {
		usuarioJSON := []byte(`{
			"nome":  "Cláudio Araújo Atualizado",
			"email": "claudio@gmail.com"
		}`)

		w := performRequest(r, http.MethodPut, "/api/v1/usuarios/1", usuarioJSON, http.Header{
			"Authorization": []string{"Bearer " + token}})

		if w.Code != http.StatusOK {
			t.Fatalf("Esperado status 200, obtido %d – body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			Data domain.Usuario `json:"data"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Falha ao decodificar resposta JSON: %v", err)
		}

		if resp.Data.ID == 0 {
			t.Fatal("ID do usuário não foi retornado")
		}

		if resp.Data.Senha != "" {
			t.Fatal("Senha do usuário não deve ser retornada")
		}
	})

	t.Run("Alterar usuário com token inválido", func(t *testing.T) {
		usuarioJSON := []byte(`{
			"nome":  "Cláudio Araújo Atualizado",
			"email": "claudio@gmail.com"
		}`)

		w := performRequest(r, http.MethodPut, "/api/v1/usuarios/1", usuarioJSON, http.Header{
			"Authorization": []string{"Bearer tokenInvalido"}})

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperado status 401, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Alterar usuário sem token", func(t *testing.T) {
		usuarioJSON := []byte(`{
			"nome":  "Cláudio Araújo Atualizado",
			"email": "claudio@gmail.com"
		}`)

		w := performRequest(r, http.MethodPut, "/api/v1/usuarios/1", usuarioJSON, nil)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperado status 401, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Alterar usuário com dados inválidos", func(t *testing.T) {
		usuarioJSON := []byte(`{
			"nome":  "",
			"email": "emailInvalido"
		}`)

		w := performRequest(r, http.MethodPut, "/api/v1/usuarios/1", usuarioJSON, http.Header{
			"Authorization": []string{"Bearer " + token}})

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

	t.Run("Alterar usuário com email já existente", func(t *testing.T) {
		novoUsuarioJSON := []byte(`{
			"nome":  "Cláudio Araújo",
			"email": "claudioteste@gmail.com",
			"senha": "123456"
		}`)

		w := performRequest(r, http.MethodPost, "/api/v1/usuarios", novoUsuarioJSON, nil)

		if w.Code != http.StatusCreated {
			t.Fatalf("Esperado status 201, obtido %d – body: %s", w.Code, w.Body.String())
		}

		usuarioJSON := []byte(`{
			"nome":  "Cláudio Araújo Atualizado",
			"email": "claudioteste@gmail.com"
		}`)

		w = performRequest(r, http.MethodPut, "/api/v1/usuarios/1", usuarioJSON, http.Header{
			"Authorization": []string{"Bearer " + token}})

		if w.Code != http.StatusConflict {
			t.Fatalf("Esperado status 409, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})
}

func TestAlterarSenhaUsuarioLogado(t *testing.T){

	t.Run("Alterar senha com dados válidos", func(t *testing.T) {
		senhaJSON := []byte(`{
			"senha_atual": "123456",
			"senha_nova": "novaSenha123"
		}`)

		w := performRequest(r, http.MethodPut, "/api/v1/usuarios/senha", senhaJSON, http.Header{
			"Authorization": []string{"Bearer " + token}})

		if w.Code != http.StatusOK {
			t.Fatalf("Esperado status 200, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Alterar senha com token inválido", func(t *testing.T) {
		senhaJSON := []byte(`{
			"senha_atual": "123456",
			"senha_nova": "novaSenha123"
		}`)

		w := performRequest(r, http.MethodPut, "/api/v1/usuarios/senha", senhaJSON, http.Header{
			"Authorization": []string{"Bearer tokenInvalido"}})

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperado status 401, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Alterar senha sem token", func(t *testing.T) {
		senhaJSON := []byte(`{
			"senha_atual": "123456",
			"senha_nova": "novaSenha123"
		}`)

		w := performRequest(r, http.MethodPut, "/api/v1/usuarios/senha", senhaJSON, nil)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("Esperado status 401, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("Alterar senha com dados inválidos", func(t *testing.T) {
		senhaJSON := []byte(`{
			"senha_atual": "123456",
			"senha_nova": "123"
		}`)

		w := performRequest(r, http.MethodPut, "/api/v1/usuarios/senha", senhaJSON, http.Header{
			"Authorization": []string{"Bearer " + token}})

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

	t.Run("Alterar senha com senha atual incorreta", func(t *testing.T) {
		senhaJSON := []byte(`{
			"senha_atual": "senhaIncorreta",
			"senha_nova": "novaSenha123"
		}`)

		w := performRequest(r, http.MethodPut, "/api/v1/usuarios/senha", senhaJSON, http.Header{
			"Authorization": []string{"Bearer " + token}})

		if w.Code != http.StatusForbidden {
			t.Fatalf("Esperado status 403, obtido %d – body: %s", w.Code, w.Body.String())
		}
	})
}
