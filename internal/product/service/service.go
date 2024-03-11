package productsvc

import (
	"context"
	"fmt"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/product"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/validation"
)

type Service struct {
	repo product.Repository
}

func New(repo product.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) Create(ctx context.Context, req model.ProductCreateRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("product.service.Create: failed to validate request: %w", err)
		return
	}

	if req.Tags == nil {
		req.Tags = []string{}
	}

	data := entity.Product{
		UserID:         req.UserID,
		Name:           req.Name,
		Price:          *req.Price,
		ImageUrl:       req.ImageUrl,
		Stock:          *req.Stock,
		Condition:      entity.ProductCondition(req.Condition),
		Tags:           req.Tags,
		IsPurchaseable: *req.IsPurchaseable,
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("product.service.Create: failed to create product: %w", err)
		return
	}

	return
}
