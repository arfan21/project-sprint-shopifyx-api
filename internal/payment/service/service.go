package paymentsvc

import (
	"context"
	"fmt"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/payment"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/product"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/validation"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Service struct {
	repo           payment.Repository
	bankAccountSvc bankaccount.Service
	productSvc     product.Service
}

func New(repo payment.Repository, bankAccountSvc bankaccount.Service, productSvc product.Service) *Service {
	return &Service{repo: repo, bankAccountSvc: bankAccountSvc, productSvc: productSvc}
}

func (s Service) Create(ctx context.Context, req model.PaymentRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("payment.service.Create: failed to validate request: %w", err)
		return
	}

	_, err = s.bankAccountSvc.GetByID(ctx, req.BankAccountID, req.UserID)
	if err != nil {
		err = fmt.Errorf("payment.service.Create: failed to get bank account: %w", err)
		return
	}

	resProduct, err := s.productSvc.GetByID(ctx, req.ProductID)
	if err != nil {
		err = fmt.Errorf("payment.service.Create: failed to get product: %w", err)
		return
	}

	if !resProduct.IsPurchaseable {
		err = fmt.Errorf("payment.service.Create: failed to purchase product: %w", constant.ErrProductNotPurchaseable)
		return
	}

	if resProduct.UserID == req.UserID {
		err = fmt.Errorf("payment.service.Create: failed to purchase own product: %w", constant.ErrCannotButOwnProduct)
		return
	}

	if resProduct.Stock < req.Quantity {
		err = fmt.Errorf("payment.service.Create: failed to purchase product: %w", constant.ErrInsufficientStock)
		return
	}

	totalPrice := resProduct.Price.Mul(decimal.NewFromInt(int64(req.Quantity)))

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("payment.service.Create: failed to generate product id: %w", err)
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
		err = fmt.Errorf("payment.service.Create: failed to begin transaction: %w", err)
		return
	}

	defer func() {
		if err != nil {
			errRb := tx.Rollback(ctx)
			if errRb != nil {
				errRb = fmt.Errorf("payment.service.Create: failed to rollback transaction: %w", errRb)
				err = errRb
				return
			}
			return
		}

		err = tx.Commit(ctx)
		if err != nil {
			err = fmt.Errorf("payment.service.Create: failed to commit transaction: %w", err)
			return
		}
	}()

	err = s.repo.WithTx(tx).Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("payment.service.Create: failed to create payment: %w", err)
		return
	}

	err = s.productSvc.WithTx(tx).ReduceStock(ctx, req.ProductID, req.Quantity)
	if err != nil {
		err = fmt.Errorf("payment.service.Create: failed to reduce stock: %w", err)
		return
	}

	return
}
