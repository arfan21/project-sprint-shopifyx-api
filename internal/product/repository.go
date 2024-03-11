package product

import (
	"context"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, data entity.Product) (err error)
	Update(ctx context.Context, data entity.Product) (err error)
	GetByID(ctx context.Context, id uuid.UUID) (product entity.Product, err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
}
