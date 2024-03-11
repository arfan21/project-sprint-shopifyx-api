package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductCreateRequest struct {
	UserID         uuid.UUID        `json:"-" validate:"required"`
	Name           string           `json:"name" validate:"required,min=5,max=60"`
	Price          *decimal.Decimal `json:"price" validate:"required"`
	ImageUrl       string           `json:"imageUrl" validate:"required,url"`
	Stock          *int             `json:"stock" validate:"required"`
	Condition      string           `json:"condition" validate:"required,oneof=new second"`
	Tags           []string         `json:"tags" validate:"omitempty,dive,min=3,max=20"`
	IsPurchaseable *bool            `json:"isPurchaseable" validate:"required"`
}
