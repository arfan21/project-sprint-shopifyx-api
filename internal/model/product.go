package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductRequest struct {
	ID             uuid.UUID        `json:"-" validate:"omitempty"`
	UserID         uuid.UUID        `json:"-" validate:"omitempty"`
	Name           string           `json:"name" validate:"required,min=5,max=60"`
	Price          *decimal.Decimal `json:"price" validate:"required"`
	ImageUrl       string           `json:"imageUrl" validate:"required,url"`
	Stock          *int             `json:"stock" validate:"required"`
	Condition      string           `json:"condition" validate:"required,oneof=new second"`
	Tags           []string         `json:"tags" validate:"omitempty,dive,min=3,max=20"`
	IsPurchaseable *bool            `json:"isPurchaseable" validate:"required"`
}

type ProductDeleteRequest struct {
	ID     uuid.UUID `json:"id" validate:"required"`
	UserID uuid.UUID `json:"-" validate:"required"`
}
