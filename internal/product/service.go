package product

import (
	"context"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
)

type Service interface {
	Create(ctx context.Context, req model.ProductRequest) (err error)
	Update(ctx context.Context, req model.ProductRequest) (err error)
	Delete(ctx context.Context, req model.ProductDeleteRequest) (err error)
	GetList(ctx context.Context, req model.ProductGetListRequest) (res []model.ProductGetResponse, total int, err error)
}
