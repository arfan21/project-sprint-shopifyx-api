package payment

import (
	"context"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	paymentrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/payment/repository"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) *paymentrepo.Repository

	Create(ctx context.Context, data entity.Payment) (err error)
}
