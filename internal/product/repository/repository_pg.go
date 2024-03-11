package productrepo

import (
	"context"
	"fmt"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	dbpostgres "github.com/arfan21/project-sprint-shopifyx-api/pkg/db/postgres"
	"github.com/google/uuid"
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
		INSERT INTO products (id, name, price, imageUrl, stock, condition, tags, isPurchaseable, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8,  $9)
	`

	_, err = r.db.Exec(ctx, query,
		data.ID,
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

func (r Repository) Update(ctx context.Context, data entity.Product) (err error) {
	query := `
		UPDATE products
		SET name = $1, price = $2, imageUrl = $3, condition = $4, tags = $5, isPurchaseable = $6
		WHERE id = $7
	`

	_, err = r.db.Exec(ctx, query,
		data.Name,
		data.Price,
		data.ImageUrl,
		data.Condition,
		data.Tags,
		data.IsPurchaseable,
		data.ID,
	)
	if err != nil {
		err = fmt.Errorf("product.repository.Update: failed to update product: %w", err)
		return
	}

	return
}

func (r Repository) GetByID(ctx context.Context, id uuid.UUID) (product entity.Product, err error) {
	query := `
		SELECT id, name, price, imageUrl, stock, condition, tags, isPurchaseable, user_id
		FROM products
		WHERE id = $1
	`

	err = r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.ImageUrl,
		&product.Stock,
		&product.Condition,
		&product.Tags,
		&product.IsPurchaseable,
		&product.UserID,
	)
	if err != nil {
		err = fmt.Errorf("product.repository.GetByID: failed to get product by id: %w", err)
		return
	}

	return
}

func (r Repository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	query := `
		DELETE FROM products
		WHERE id = $1
	`

	_, err = r.db.Exec(ctx, query, id)
	if err != nil {
		err = fmt.Errorf("product.repository.Delete: failed to delete product: %w", err)
		return
	}

	return
}
