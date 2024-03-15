package productsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/product"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type Service struct {
	repo           product.Repository
	bankAccountSvc bankaccount.Service
}

func New(repo product.Repository, bankAccountSvc bankaccount.Service) *Service {
	return &Service{repo: repo, bankAccountSvc: bankAccountSvc}
}

func (s Service) WithTx(tx pgx.Tx) product.Service {
	s.repo = s.repo.WithTx(tx)
	return &s
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

func (s Service) GetList(ctx context.Context, req model.ProductGetListRequest) (res []model.ProductDetailResponse, total int, err error) {
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

	res = make([]model.ProductDetailResponse, len(resDB))
	productIds := make([]uuid.UUID, len(resDB))
	for i, v := range resDB {
		res[i] = model.ProductDetailResponse{
			ProductID:     v.ID,
			Name:          v.Name,
			Price:         v.Price.InexactFloat64(),
			ImageUrl:      v.ImageUrl,
			Stock:         v.Stock,
			Condition:     string(v.Condition),
			Tags:          v.Tags,
			IsPurchasable: v.IsPurchaseable,
		}

		productIds[i] = v.ID
	}

	// errg, ctxg := errgroup.WithContext(ctx)

	// errg.Go(func() error {
	// 	total, err = s.repo.GetTotal(ctxg, req)
	// 	if err != nil {
	// 		err = fmt.Errorf("product.service.GetList: failed to get total product: %w", err)
	// 		return err
	// 	}

	// 	return nil
	// })

	// errg.Go(func() error {
	// 	purchaseCountsMap, err := s.repo.GetPurchaseCountByProductIds(ctxg, productIds)
	// 	if err != nil {
	// 		err = fmt.Errorf("product.service.GetList: failed to get purchase count by product ids: %w", err)
	// 		return err
	// 	}

	// 	for i, v := range res {
	// 		if purchaseCount, ok := purchaseCountsMap[v.ProductID]; ok {
	// 			res[i].PurchaseCount = purchaseCount
	// 		}
	// 	}

	// 	return nil
	// })

	// err = errg.Wait()
	// if err != nil {
	// 	err = fmt.Errorf("product.service.GetList: failed to wait errgroup: %w", err)
	// 	return
	// }

	total, err = s.repo.GetTotal(ctx, req)
	if err != nil {
		err = fmt.Errorf("product.service.GetList: failed to get total product: %w", err)
		return
	}

	purchaseCountsMap, err := s.repo.GetPurchaseCountByProductIds(ctx, productIds)
	if err != nil {
		err = fmt.Errorf("product.service.GetList: failed to get purchase count by product ids: %w", err)
		return
	}

	for i, v := range res {
		if purchaseCount, ok := purchaseCountsMap[v.ProductID]; ok {
			res[i].PurchaseCount = purchaseCount
		}
	}

	return
}

func (s Service) GetDetailByID(ctx context.Context, id uuid.UUID) (res model.ProductDetailResponse, err error) {
	resDB, err := s.repo.GetDetailByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrProductNotFound
		}
		err = fmt.Errorf("product.service.GetDetailByID: failed to get product by id: %w", err)
		return
	}

	res = model.ProductDetailResponse{
		ProductID:     resDB.ID,
		Name:          resDB.Name,
		Price:         resDB.Price.InexactFloat64(),
		ImageUrl:      resDB.ImageUrl,
		Stock:         resDB.Stock,
		Condition:     string(resDB.Condition),
		Tags:          resDB.Tags,
		IsPurchasable: resDB.IsPurchaseable,
	}

	purchaseCountsMap, err := s.repo.GetPurchaseCountByProductIds(ctx, []uuid.UUID{resDB.ID})
	if err != nil {
		err = fmt.Errorf("product.service.GetList: failed to get purchase count by product ids: %w", err)
		return
	}

	if purchaseCount, ok := purchaseCountsMap[resDB.ID]; ok {
		res.PurchaseCount = purchaseCount
	}

	purchaseCountSeller, err := s.repo.GetPurchaseCountBySeller(ctx, resDB.UserID)
	if err != nil {
		err = fmt.Errorf("product.service.GetDetailByID: failed to get purchase count by seller: %w", err)
		return
	}

	sellerBankAccounts, err := s.bankAccountSvc.GetListByUserID(ctx, resDB.UserID)
	if err != nil {
		err = fmt.Errorf("product.service.GetDetailByID: failed to get bank accounts by user id: %w", err)
		return
	}

	res.Seller = model.ProductDetailSellerResponse{
		ProductSoldTotal: purchaseCountSeller,
		Name:             resDB.Seller.Name,
		BankAccounts:     sellerBankAccounts,
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
			errRb := tx.Rollback(ctx)
			if errRb != nil {
				errRb = fmt.Errorf("product.service.Payment: failed to rollback transaction: %w", errRb)
				err = errRb
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

func (s Service) ReduceStock(ctx context.Context, id uuid.UUID, qty int) (err error) {
	return s.repo.ReduceStock(ctx, id, qty)
}

func (s Service) GetByID(ctx context.Context, id uuid.UUID) (res model.ProductGetResponse, err error) {
	resDB, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrProductNotFound
		}
		err = fmt.Errorf("product.service.Update: failed to get product by id: %w", err)
		return
	}

	res = model.ProductGetResponse{
		ProductID:     resDB.ID,
		Name:          resDB.Name,
		Price:         resDB.Price,
		ImageUrl:      resDB.ImageUrl,
		Stock:         resDB.Stock,
		Condition:     string(resDB.Condition),
		Tags:          resDB.Tags,
		IsPurchasable: resDB.IsPurchaseable,
		UserID:        resDB.UserID,
	}

	return
}

func (s Service) Payment(ctx context.Context, req model.PaymentRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("product.service.Payment: failed to validate request: %w", err)
		return
	}

	_, err = s.bankAccountSvc.GetByID(ctx, req.BankAccountID, req.UserID)
	if err != nil {
		err = fmt.Errorf("product.service.Payment: failed to get bank account: %w", err)
		return
	}

	resProduct, err := s.GetByID(ctx, req.ProductID)
	if err != nil {
		err = fmt.Errorf("product.service.Payment: failed to get product: %w", err)
		return
	}

	if !resProduct.IsPurchasable {
		err = fmt.Errorf("product.service.Payment: failed to purchase product: %w", constant.ErrProductNotPurchaseable)
		return
	}

	if resProduct.UserID == req.UserID {
		err = fmt.Errorf("product.service.Payment: failed to purchase own product: %w", constant.ErrCannotButOwnProduct)
		return
	}

	if resProduct.Stock < req.Quantity {
		err = fmt.Errorf("product.service.Payment: failed to purchase product: %w", constant.ErrInsufficientStock)
		return
	}

	totalPrice := resProduct.Price.Mul(decimal.NewFromInt(int64(req.Quantity)))

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("product.service.Payment: failed to generate product id: %w", err)
		return
	}

	data := entity.Payment{
		ID:                   id,
		UserID:               req.UserID,
		ProductID:            req.ProductID,
		BankAccountID:        req.BankAccountID,
		Quantity:             req.Quantity,
		TotalPrice:           totalPrice,
		PaymentProofImageURL: req.PaymentProofImageURL,
	}

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		err = fmt.Errorf("product.service.Payment: failed to begin transaction: %w", err)
		return
	}

	defer func() {
		if err != nil {
			errRb := tx.Rollback(ctx)
			if errRb != nil {
				err = fmt.Errorf("product.service.Payment: failed to rollback transaction: %w", errRb)
				return
			}
			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			err = fmt.Errorf("product.service.Payment: failed to commit transaction: %w", err)
			return
		}
	}()

	err = s.repo.WithTx(tx).Payment(ctx, data)
	if err != nil {
		err = fmt.Errorf("product.service.Payment: failed to create payment: %w", err)
		return
	}

	err = s.repo.WithTx(tx).ReduceStock(ctx, req.ProductID, req.Quantity)
	if err != nil {
		err = fmt.Errorf("product.service.Payment: failed to reduce stock: %w", err)
		return
	}

	return
}
