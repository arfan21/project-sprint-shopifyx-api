package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Payment struct {
	ID                   uuid.UUID       `json:"id"`
	UserID               uuid.UUID       `json:"userId"`
	ProductID            uuid.UUID       `json:"productId"`
	BankAccountID        uuid.UUID       `json:"bankAccountId"`
	PaymentProofImageURL string          `json:"paymentProofImageUrl"`
	Quantity             int             `json:"quantity"`
	TotalPrice           decimal.Decimal `json:"totalPrice"`
	CreatedAt            time.Time       `json:"createdAt"`
	UpdatedAt            time.Time       `json:"updatedAt"`
}

func (Payment) TableName() string {
	return "payments"
}
