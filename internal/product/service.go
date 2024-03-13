package product

import (
	"context"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	WithTx(tx pgx.Tx) Service

	Create(ctx context.Context, req model.ProductRequest) (err error)
	Update(ctx context.Context, req model.ProductRequest) (err error)
	Delete(ctx context.Context, req model.ProductDeleteRequest) (err error)
	GetList(ctx context.Context, req model.ProductGetListRequest) (res []model.ProductDetailResponse, total int, err error)
	GetDetailByID(ctx context.Context, id uuid.UUID) (res model.ProductDetailResponse, err error)
	UpdateStock(ctx context.Context, req model.ProductUpdateStockRequest) (err error)
	ReduceStock(ctx context.Context, id uuid.UUID, qty int) (err error)
	GetByID(ctx context.Context, id uuid.UUID) (res model.ProductGetResponse, err error)
}
