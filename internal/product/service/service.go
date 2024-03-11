package productsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/product"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo product.Repository
}

func New(repo product.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) Create(ctx context.Context, req model.ProductRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("product.service.Create: failed to validate request: %w", err)
		return
	}

	if req.Tags == nil {
		req.Tags = []string{}
	}

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("product.service.Create: failed to generate product id: %w", err)
		return
	}

	data := entity.Product{
		ID:             id,
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

func (s Service) Update(ctx context.Context, req model.ProductRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("product.service.Update: failed to validate request: %w", err)
		return
	}

	err = s.validateProduct(ctx, req.ID, req.UserID)
	if err != nil {
		return
	}

	data := entity.Product{
		ID:             req.ID,
		Name:           req.Name,
		Price:          *req.Price,
		ImageUrl:       req.ImageUrl,
		Stock:          *req.Stock,
		Condition:      entity.ProductCondition(req.Condition),
		Tags:           req.Tags,
		IsPurchaseable: *req.IsPurchaseable,
	}

	err = s.repo.Update(ctx, data)
	if err != nil {
		err = fmt.Errorf("product.service.Update: failed to update product: %w", err)
		return
	}

	return
}

func (s Service) Delete(ctx context.Context, req model.ProductDeleteRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("product.service.Update: failed to validate request: %w", err)
		return
	}

	err = s.validateProduct(ctx, req.ID, req.UserID)
	if err != nil {
		return
	}

	err = s.repo.Delete(ctx, req.ID)
	if err != nil {
		err = fmt.Errorf("product.service.Delete: failed to delete product: %w", err)
		return
	}

	return
}

func (s Service) validateProduct(ctx context.Context, id, userID uuid.UUID) (err error) {
	resDB, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrProductNotFound
		}
		err = fmt.Errorf("product.service.Update: failed to get product by id: %w", err)
		return
	}

	if resDB.UserID != userID {
		err = fmt.Errorf("product.service.Update: user id not match, %w", constant.ErrAccessForbidden)
		return
	}

	return
}
