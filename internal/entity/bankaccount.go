package entity

import (
	"time"

	"github.com/google/uuid"
)

type BankAccount struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"userId"`
	BankName      string    `json:"bankName"`
	AccountNumber string    `json:"accountNumber"`
	AccountHolder string    `json:"accountHolder"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (BankAccount) TableName() string {
	return "bank_accounts"
}
