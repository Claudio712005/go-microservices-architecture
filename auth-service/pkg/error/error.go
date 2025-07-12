package error

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Kind categoriza o erro.
type Kind int

const (
	_ Kind = iota
	KindInternal
	KindNotFound
	KindValidation
	KindUnauthorized
	KindForbidden
	KindConflict
)

// AppError é o tipo de erro usado na aplicação.
type AppError struct {
	Kind    Kind   `json:"-"`          // não expõe Kind
	Code    string `json:"code"`       // ex.: USER_NOT_FOUND
	Err     error  `json:"-"`          // erro original
	Message string `json:"message"`    // msg “safe” para cliente
	Status  int    `json:"-"`          // sugestão de HTTP status
}

// Error implementa a interface error.
func (e *AppError) Error() string  { return e.Err.Error() }
// Unwrap é usado para acessar o erro original.
func (e *AppError) Unwrap() error  { return e.Err }
// String é uma representação legível do erro, incluindo Kind e Code.
func (e *AppError) String() string { return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err) }

// Serializa somente Code e Message.
func (e *AppError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{Code: e.Code, Message: e.Message})
}

// New cria um AppError do zero.
func New(kind Kind, code, msg string, status int, err error) *AppError {
	if err == nil {
		err = errors.New(msg)
	}
	return &AppError{kind, code, err, msg, status}
}

// NotFound cria um AppError do tipo NotFound.
func NotFound(code, msg string, err error) *AppError {
	return New(KindNotFound, code, msg, http.StatusNotFound, err)
}
// Validation cria um AppError do tipo Validation.
func Validation(code, msg string, err error) *AppError {
	return New(KindValidation, code, msg, http.StatusBadRequest, err)
}
// Unauthorized cria um AppError do tipo Unauthorized.
func Unauthorized(code, msg string, err error) *AppError {
	return New(KindUnauthorized, code, msg, http.StatusUnauthorized, err)
}
// Forbidden cria um AppError do tipo Forbidden.
func Forbidden(code, msg string, err error) *AppError {
	return New(KindForbidden, code, msg, http.StatusForbidden, err)
}
// Internal cria um AppError do tipo Internal.
func Internal(code, msg string, err error) *AppError {
	return New(KindInternal, code, msg, http.StatusInternalServerError, err)
}
// Conflict cria um AppError do tipo Conflict.
func Conflict(code, msg string, err error) *AppError {
	return New(KindConflict, code, msg, http.StatusConflict, err)
}

// Wrap mantém a stack de erros, mas os “empacota” em AppError.
func Wrap(err error, kind Kind, code, msg string, status int) *AppError {
	return New(kind, code, msg, status, err)
}
