package product

import (
	"context"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, data entity.Product) (err error)
}
