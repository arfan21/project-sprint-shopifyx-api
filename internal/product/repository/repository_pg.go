package productrepo

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

func (r Repository) Create(ctx context.Context, data entity.Product) (err error) {
	query := `
		INSERT INTO products (name, price, imageUrl, stock, condition, tags, isPurchaseable, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = r.db.Exec(ctx, query,
		data.Name,
		data.Price,
		data.ImageUrl,
		data.Stock,
		data.Condition,
		data.Tags,
		data.IsPurchaseable,
		data.UserID,
	)
	if err != nil {
		err = fmt.Errorf("product.repository.Create: failed to create product: %w", err)
		return
	}

	return
}
