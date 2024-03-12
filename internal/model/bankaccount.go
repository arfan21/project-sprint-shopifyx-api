package model

import "github.com/google/uuid"

type BankAccountRequest struct {
	BankName          string    `json:"bankName" validate:"required"`
	BankAccountName   string    `json:"bankAccountName" validate:"required"`
	BankAccountNumber string    `json:"bankAccountNumber" validate:"required"`
	BankAccountID     uuid.UUID `json:"bankAccountId"`
	UserID            uuid.UUID `json:"-"`
}
