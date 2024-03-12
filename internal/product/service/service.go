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

func (s Service) GetList(ctx context.Context, req model.ProductGetListRequest) (res []model.ProductGetResponse, total int, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("product.service.GetList: failed to validate request: %w", err)
		return
	}

	resDB, err := s.repo.GetList(ctx, req)
	if err != nil {
		err = fmt.Errorf("product.service.GetList: failed to get product list: %w", err)
		return
	}

	res = make([]model.ProductGetResponse, len(resDB))

	for i, v := range resDB {
		res[i] = model.ProductGetResponse{
			ProductID:      v.ID,
			Name:           v.Name,
			Price:          v.Price.InexactFloat64(),
			ImageUrl:       v.ImageUrl,
			Stock:          v.Stock,
			Condition:      string(v.Condition),
			Tags:           v.Tags,
			IsPurchaseable: v.IsPurchaseable,
		}
	}

	total, err = s.repo.GetTotal(ctx, req)
	if err != nil {
		err = fmt.Errorf("product.service.GetList: failed to get total product: %w", err)
		return
	}

	// TODO: add get purchase count each product

	return
}

func (s Service) GetDetailByID(ctx context.Context, id uuid.UUID) (res model.ProductGetResponse, err error) {
	resDB, err := s.repo.GetDetailByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrProductNotFound
		}
		err = fmt.Errorf("product.service.GetDetailByID: failed to get product by id: %w", err)
		return
	}

	res = model.ProductGetResponse{
		ProductID:      resDB.ID,
		Name:           resDB.Name,
		Price:          resDB.Price.InexactFloat64(),
		ImageUrl:       resDB.ImageUrl,
		Stock:          resDB.Stock,
		Condition:      string(resDB.Condition),
		Tags:           resDB.Tags,
		IsPurchaseable: resDB.IsPurchaseable,
	}

	return
}

func (s Service) UpdateStock(ctx context.Context, req model.ProductUpdateStockRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("product.service.UpdateStock: failed to validate request: %w", err)
		return
	}

	err = s.validateProduct(ctx, req.ID, req.UserID)
	if err != nil {
		return
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("product.service.UpdateStock: failed to begin transaction: %w", err)
		return
	}

	defer func() {
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				err = fmt.Errorf("product.service.UpdateStock: failed to rollback transaction: %w", err)
				return
			}
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			err = fmt.Errorf("product.service.UpdateStock: failed to commit transaction: %w", err)
			return
		}
	}()

	_, err = s.repo.WithTx(tx).GetStockByIDForUpdate(ctx, req.ID)
	if err != nil {
		err = fmt.Errorf("product.service.UpdateStock: failed to get stock for update: %w", err)
		return
	}

	err = s.repo.WithTx(tx).UpdateStock(ctx, req.ID, req.Stock)
	if err != nil {
		err = fmt.Errorf("product.service.UpdateStock: failed to update stock: %w", err)
		return
	}

	return
}
