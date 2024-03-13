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

type ProductGetListRequest struct {
	UserOnly       bool       `json:"userOnly" query:"userOnly"`
	UserID         *uuid.UUID `json:"-" validate:"required_if=UserOnly true"`
	Limit          int        `json:"limit" query:"limit" validate:"required_unless=Offset 0"`
	Offset         int        `json:"offset" query:"offset"`
	Condition      string     `json:"condition" query:"condition" validate:"omitempty,oneof=new second"`
	Tags           []string   `json:"tags" query:"tags" validate:"omitempty,dive,min=3,max=20"`
	ShowEmptyStock bool       `json:"showEmptyStock" query:"showEmptyStock"`
	MaxPrice       float64    `json:"maxPrice" query:"maxPrice" validate:"required_unless=MinPrice 0,omitempty,gtfield=MinPrice"`
	MinPrice       float64    `json:"minPrice" query:"minPrice" validate:"required_unless=MaxPrice 0,omitempty,ltfield=MaxPrice"`
	SortBy         string     `json:"sortBy" query:"sortBy" validate:"omitempty,oneof=price date"`
	OrderBy        string     `json:"orderBy" query:"orderBy" validate:"omitempty,oneof=asc desc"`
	Search         string     `json:"search" query:"search" validate:"omitempty,min=3,max=60"`
	DisablePaging  bool       `json:"-"`
	DisableOrder   bool       `json:"-"`
}

type ProductDetailResponse struct {
	ProductID      uuid.UUID `json:"productId"`
	Name           string    `json:"name"`
	Price          float64   `json:"price"`
	ImageUrl       string    `json:"imageUrl"`
	Stock          int       `json:"stock"`
	Condition      string    `json:"condition"`
	Tags           []string  `json:"tags"`
	IsPurchaseable bool      `json:"isPurchaseable"`
	PurchaseCount  int       `json:"purchaseCount"`
}

// ProductGetResponse is model for internal use
type ProductGetResponse struct {
	ProductID      uuid.UUID       `json:"productId"`
	UserID         uuid.UUID       `json:"userId"`
	Name           string          `json:"name"`
	Price          decimal.Decimal `json:"price"`
	ImageUrl       string          `json:"imageUrl"`
	Stock          int             `json:"stock"`
	Condition      string          `json:"condition"`
	Tags           []string        `json:"tags"`
	IsPurchaseable bool            `json:"isPurchaseable"`
}

type ProductUpdateStockRequest struct {
	ID     uuid.UUID `json:"-" validate:"required"`
	UserID uuid.UUID `json:"-" validate:"required"`
	Stock  int       `json:"stock" validate:"required"`
}
