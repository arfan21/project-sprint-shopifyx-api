package userrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/arfan21/shopifyx-api/internal/entity"
	"github.com/arfan21/shopifyx-api/pkg/constant"
	dbpostgres "github.com/arfan21/shopifyx-api/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db dbpostgres.Queryer
}

func New(db dbpostgres.Queryer) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) Begin(ctx context.Context) (tx pgx.Tx, err error) {
	return r.db.Begin(ctx)
}

func (r Repository) WithTx(tx pgx.Tx) *Repository {
	r.db = tx
	return &r
}

func (r Repository) Create(ctx context.Context, data entity.User) (err error) {
	query := `
		INSERT INTO users (username, name, password)
		VALUES ($1, $2, $3)
	`

	_, err = r.db.Exec(ctx, query, data.Username, data.Name, data.Password)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == constant.ErrSQLUniqueViolation {
				err = constant.ErrUsernameAlreadyRegistered
			}
		}

		err = fmt.Errorf("user.repository.Create: failed to create user: %w", err)
		return
	}

	return
}

func (r Repository) GetByUsername(ctx context.Context, username string) (data entity.User, err error) {
	query := `
		SELECT id, username, name, password
		FROM users
		WHERE username = $1
	`

	err = r.db.QueryRow(ctx, query, username).Scan(
		&data.ID,
		&data.Username,
		&data.Name,
		&data.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrUsernameOrPasswordInvalid
		}

		err = fmt.Errorf("user.repository.GetByEmail: failed to get user by email: %w", err)

		return
	}

	return
}
