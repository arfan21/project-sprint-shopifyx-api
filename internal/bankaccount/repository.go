package bankaccount

import (
	"context"

	bankaccountrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/repository"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) *bankaccountrepo.Repository

	Create(ctx context.Context, data entity.BankAccount) (err error)
	Update(ctx context.Context, data entity.BankAccount) (err error)
	GetByID(ctx context.Context, id, userId uuid.UUID) (bankAccount entity.BankAccount, err error)
	Delete(ctx context.Context, id uuid.UUID) (err error)
	GetListByUserID(ctx context.Context, userId uuid.UUID) (res []entity.BankAccount, err error)
}
