package bankaccount

import (
	"context"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, req model.BankAccountRequest) (err error)
	Update(ctx context.Context, req model.BankAccountRequest) (err error)
	Delete(ctx context.Context, id, userId uuid.UUID) (err error)
	GetListByUserID(ctx context.Context, userId uuid.UUID) (res []model.BankAccountResponse, err error)
	GetByID(ctx context.Context, id, userId uuid.UUID) (res model.BankAccountResponse, err error)
}
