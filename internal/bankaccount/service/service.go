package bankaccountsvc

import (
	"context"
	"fmt"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/validation"
	"github.com/google/uuid"
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
