package bankaccount

import (
	"context"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
)

type Service interface {
	Create(ctx context.Context, req model.BankAccountRequest) (err error)
	Update(ctx context.Context, req model.BankAccountRequest) (err error)
}
