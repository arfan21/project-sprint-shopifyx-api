package constant

import (
	"errors"
	"net/http"
)

const (
	ErrSQLUniqueViolation = "23505"
)

var (
	ErrUsernameAlreadyRegistered = &ErrWithCode{HTTPStatusCode: http.StatusConflict, Message: "username already registered"}
	ErrUsernameOrPasswordInvalid = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "username or password invalid"}
	ErrUserNotFound              = &ErrWithCode{HTTPStatusCode: http.StatusNotFound, Message: "user not found"}
	ErrInvalidUUID               = errors.New("invalid uuid length or format")
	ErrAccessForbidden           = &ErrWithCode{HTTPStatusCode: http.StatusForbidden, Message: "access forbidden"}
	ErrProductNotFound           = &ErrWithCode{HTTPStatusCode: http.StatusNotFound, Message: "product not found"}
	ErrBankAccountNotFound       = &ErrWithCode{HTTPStatusCode: http.StatusNotFound, Message: "bank account not found"}
	ErrInsufficientStock         = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "insufficient stock"}
	ErrCannotButOwnProduct       = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "cannot buy own product"}
	ErrProductNotPurchaseable    = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "product not purchaseable"}
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
