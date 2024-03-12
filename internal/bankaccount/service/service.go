package bankaccountsvc

import (
	"context"
	"errors"
	"fmt"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo bankaccount.Repository
}

func New(repo bankaccount.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) Create(ctx context.Context, req model.BankAccountRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("bankaccount.service.Create: failed to validate request: %w", err)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		err = fmt.Errorf("bankaccount.service.Create: failed to generate product id: %w", err)
		return
	}

	data := entity.BankAccount{
		ID:            id,
		BankName:      req.BankName,
		AccountNumber: req.BankAccountNumber,
		AccountHolder: req.BankAccountName,
		UserID:        req.UserID,
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("bankaccount.service.Create: failed to create bank account: %w", err)
		return
	}

	return

}

func (s Service) Update(ctx context.Context, req model.BankAccountRequest) (err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("bankaccount.service.Update: failed to validate request: %w", err)
		return
	}

	_, err = s.repo.GetByID(ctx, req.BankAccountID, req.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrBankAccountNotFound
		}
		err = fmt.Errorf("bankaccount.service.Update: failed to get bank account: %w", err)
		return
	}

	data := entity.BankAccount{
		ID:            req.BankAccountID,
		BankName:      req.BankName,
		AccountNumber: req.BankAccountNumber,
		AccountHolder: req.BankAccountName,
		UserID:        req.UserID,
	}

	err = s.repo.Update(ctx, data)
	if err != nil {
		err = fmt.Errorf("bankaccount.service.Update: failed to update bank account: %w", err)
		return
	}

	return
}

func (s Service) Delete(ctx context.Context, id, userId uuid.UUID) (err error) {
	_, err = s.repo.GetByID(ctx, id, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrBankAccountNotFound
		}
		err = fmt.Errorf("bankaccount.service.Delete: failed to get bank account: %w", err)
		return
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		err = fmt.Errorf("bankaccount.service.Delete: failed to delete bank account: %w", err)
		return
	}

	return
}

func (s Service) GetListByUserID(ctx context.Context, userId uuid.UUID) (res []model.BankAccountResponse, err error) {
	resDB, err := s.repo.GetListByUserID(ctx, userId)
	if err != nil {
		err = fmt.Errorf("bankaccount.service.GetListByUserID: failed to get bank account by user id: %w", err)
		return
	}

	res = make([]model.BankAccountResponse, len(resDB))

	for i, v := range resDB {
		res[i] = model.BankAccountResponse{
			BankAccountID:     v.ID,
			BankName:          v.BankName,
			BankAccountNumber: v.AccountNumber,
			BankAccountName:   v.AccountHolder,
		}
	}

	return
}
