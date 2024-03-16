package model

import "github.com/google/uuid"

type BankAccountRequest struct {
	BankName          string    `json:"bankName" validate:"required"`
	BankAccountName   string    `json:"bankAccountName" validate:"required"`
	BankAccountNumber string    `json:"bankAccountNumber" validate:"required,stringuint"`
	BankAccountID     uuid.UUID `json:"bankAccountId"`
	UserID            uuid.UUID `json:"-"`
}

type BankAccountResponse struct {
	BankName          string    `json:"bankName" `
	BankAccountName   string    `json:"bankAccountName" `
	BankAccountNumber string    `json:"bankAccountNumber"`
	BankAccountID     uuid.UUID `json:"bankAccountId"`
}
