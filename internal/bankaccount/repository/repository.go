package bankaccountrepo

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

func (r Repository) Create(ctx context.Context, data entity.BankAccount) (err error) {
	query := `
		INSERT INTO bank_accounts (id, accountNumber, accountHolder, bankName, userId)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.Exec(ctx, query,
		data.ID,
		data.AccountNumber,
		data.AccountHolder,
		data.BankName,
		data.UserID,
	)
	if err != nil {
		err = fmt.Errorf("bankaccount.repository.Create: failed to create bank account: %w", err)
		return
	}

	return
}

func (r Repository) Update(ctx context.Context, data entity.BankAccount) (err error) {
	query := `
		UPDATE bank_accounts
		SET accountNumber = $1, accountHolder = $2, bankName = $3
		WHERE id = $4
	`

	_, err = r.db.Exec(ctx, query,
		data.AccountNumber,
		data.AccountHolder,
		data.BankName,
		data.ID,
	)
	if err != nil {
		err = fmt.Errorf("bankaccount.repository.Update: failed to update bank account: %w", err)
		return
	}

	return
}

func (r Repository) Delete(ctx context.Context, id uuid.UUID) (err error) {
	query := `
		DELETE FROM bank_accounts
		WHERE id = $1
	`

	_, err = r.db.Exec(ctx, query, id)
	if err != nil {
		err = fmt.Errorf("bankaccount.repository.Delete: failed to delete bank account: %w", err)
		return
	}

	return
}

func (r Repository) GetByID(ctx context.Context, id, userId uuid.UUID) (bankAccount entity.BankAccount, err error) {
	query := `
		SELECT id, accountNumber, accountHolder, bankName, userId
		FROM bank_accounts
		WHERE id = $1 AND userId = $2
	`

	err = r.db.QueryRow(ctx, query, id, userId).Scan(
		&bankAccount.ID,
		&bankAccount.AccountNumber,
		&bankAccount.AccountHolder,
		&bankAccount.BankName,
		&bankAccount.UserID,
	)
	if err != nil {
		err = fmt.Errorf("bankaccount.repository.GetByID: failed to get bank account: %w", err)
		return
	}

	return
}

func (r Repository) GetListByUserID(ctx context.Context, userID uuid.UUID) (res []entity.BankAccount, err error) {
	query := `
		SELECT id, accountNumber, accountHolder, bankName, userId
		FROM bank_accounts
		WHERE userId = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		err = fmt.Errorf("bankaccount.repository.GetListByUserID: failed to get bank account list: %w", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bankAccount entity.BankAccount
		err = rows.Scan(
			&bankAccount.ID,
			&bankAccount.AccountNumber,
			&bankAccount.AccountHolder,
			&bankAccount.BankName,
			&bankAccount.UserID,
		)
		if err != nil {
			err = fmt.Errorf("bankaccount.repository.GetListByUserID: failed to scan bank account: %w", err)
			return
		}

		res = append(res, bankAccount)
	}

	return
}
