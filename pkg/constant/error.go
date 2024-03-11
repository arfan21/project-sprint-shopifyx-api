package constant

import (
	"net/http"
)

const (
	ErrSQLUniqueViolation = "23505"
)

var (
	ErrUsernameAlreadyRegistered = &ErrWithCode{HTTPStatusCode: http.StatusConflict, Message: "username already registered"}
	ErrUsernameOrPasswordInvalid = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "username or password invalid"}
)

type ErrWithCode struct {
	HTTPStatusCode int
	Message        string
}

func (e *ErrWithCode) Error() string {
	return e.Message
}

type ErrValidation struct {
	Message string
}

func (e *ErrValidation) Error() string {
	return e.Message
}
