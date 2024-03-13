package product

import (
	"context"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	productrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/product/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) *productrepo.Repository

	Create(ctx context.Context, data entity.Product) (err error)
	Update(ctx context.Context, data entity.Product) (err error)
	GetByID(ctx context.Context, id uuid.UUID) (product entity.Product, err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetList(ctx context.Context, filter model.ProductGetListRequest) (res []entity.Product, err error)
	GetTotal(ctx context.Context, filter model.ProductGetListRequest) (total int, err error)
	GetDetailByID(ctx context.Context, id uuid.UUID) (product entity.Product, err error)
	GetStockByIDForUpdate(ctx context.Context, id uuid.UUID) (stock int, err error)
	UpdateStock(ctx context.Context, id uuid.UUID, stock int) (err error)
	ReduceStock(ctx context.Context, id uuid.UUID, qty int) (err error)
}
