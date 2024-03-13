package paymentrepo

import (
	"context"
	"fmt"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	dbpostgres "github.com/arfan21/project-sprint-shopifyx-api/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db dbpostgres.Queryer
}

func New(db dbpostgres.Queryer) *Repository {
	return &Repository{db: db}
}

func (r Repository) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	return r.db.Begin(ctx)
}

func (r Repository) WithTx(tx pgx.Tx) *Repository {
	r.db = tx
	return &r
}

func (r Repository) Create(ctx context.Context, data entity.Payment) (err error) {
	query := `
		INSERT INTO payments (id, userId, productId, bankAccountId, paymentProofImageUrl, quantity, totalPrice)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.Exec(ctx, query,
		data.ID,
		data.UserID,
		data.ProductID,
		data.BankAccountID,
		data.PaymentProofImageURL,
		data.Quantity,
		data.TotalPrice,
	)
	if err != nil {
		err = fmt.Errorf("payment.repository.Create: failed to create payment: %w", err)
		return
	}

	return
}
