package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductCondition string

const (
	ProductConditionNew  ProductCondition = "new"
	ProductConditionUsed ProductCondition = "used"
)

type Product struct {
	ID             uuid.UUID        `json:"id"`
	UserID         uuid.UUID        `json:"userId"`
	Name           string           `json:"name"`
	Price          decimal.Decimal  `json:"price"`
	ImageUrl       string           `json:"imageUrl"`
	Stock          int              `json:"stock"`
	Condition      ProductCondition `json:"condition"`
	Tags           []string         `json:"tags"`
	IsPurchaseable bool             `json:"isPurchaseable"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
	Seller         User             `json:"seller"`
}

func (Product) TableName() string {
	return "products"
}
