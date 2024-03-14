package model

import "github.com/google/uuid"

type PaymentRequest struct {
	ProductID            uuid.UUID `json:"productId" validate:"required"`
	BankAccountID        uuid.UUID `json:"bankAccountId" validate:"required"`
	UserID               uuid.UUID `json:"-" validate:"required"`
	PaymentProofImageURL string    `json:"paymentProofImageUrl" validate:"required,url"`
	Quantity             int       `json:"quantity" validate:"required,min=1"`
}
