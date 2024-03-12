package bankaccountsvc

import (
	"context"
	"testing"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount"
	bankaccountrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/repository"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

var (
	bankaccountSvc bankaccount.Service
	pgxMock        pgxmock.PgxPoolIface
	repoPG         bankaccount.Repository
)

func initDep(t *testing.T) {
	if pgxMock == nil {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)

		pgxMock = mockPool
	}

	if repoPG == nil {
		repoPG = bankaccountrepo.New(pgxMock)
	}

	if bankaccountSvc == nil {
		bankaccountSvc = New(repoPG)
	}
}

func TestCreate(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.BankAccountRequest{
			BankAccountID:     uuid.New(),
			UserID:            uuid.New(),
			BankName:          "BCA",
			BankAccountNumber: "1234567890",
			BankAccountName:   "Test",
		}

		pgxMock.ExpectExec("INSERT INTO bank_accounts (.+)").
			WithArgs(pgxmock.AnyArg(), req.BankAccountNumber, req.BankAccountName, req.BankName, req.UserID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := bankaccountSvc.Create(context.Background(), req)
		assert.NoError(t, err)
	})

	t.Run("failed, validation error", func(t *testing.T) {
		req := model.BankAccountRequest{
			BankAccountID:   uuid.New(),
			UserID:          uuid.New(),
			BankName:        "BCA",
			BankAccountName: "",
		}

		err := bankaccountSvc.Create(context.Background(), req)
		assert.Error(t, err)
		var validationErr *constant.ErrValidation
		assert.ErrorAs(t, err, &validationErr)
	})

}

func TestUpdate(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.BankAccountRequest{
			BankAccountID:     uuid.New(),
			UserID:            uuid.New(),
			BankName:          "BCA",
			BankAccountNumber: "1234567890",
			BankAccountName:   "Test",
		}
		pgxMock.ExpectQuery("SELECT (.+)").
			WithArgs(req.BankAccountID, req.UserID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "accountNumber", "accountHolder", "bankName", "userId"}).
				AddRow(req.BankAccountID, req.BankAccountNumber, req.BankAccountName, req.BankName, req.UserID))
		pgxMock.ExpectExec("UPDATE bank_accounts (.+)").
			WithArgs(req.BankAccountNumber, req.BankAccountName, req.BankName, req.BankAccountID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := bankaccountSvc.Update(context.Background(), req)
		assert.NoError(t, err)
	})

	t.Run("failed, validation error", func(t *testing.T) {
		req := model.BankAccountRequest{
			BankAccountID:     uuid.New(),
			UserID:            uuid.New(),
			BankName:          "BCA",
			BankAccountNumber: "1234567890",
			BankAccountName:   "",
		}

		err := bankaccountSvc.Update(context.Background(), req)
		assert.Error(t, err)
		var validationErr *constant.ErrValidation
		assert.ErrorAs(t, err, &validationErr)
	})
}

func TestDelete(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		userId := uuid.New()
		pgxMock.ExpectQuery("SELECT (.+)").
			WithArgs(id, userId).
			WillReturnRows(pgxmock.NewRows([]string{"id", "accountNumber", "accountHolder", "bankName", "userId"}).
				AddRow(id, "1234567890", "Test", "BCA", userId))
		pgxMock.ExpectExec("DELETE FROM bank_accounts (.+)").
			WithArgs(id).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := bankaccountSvc.Delete(context.Background(), id, userId)
		assert.NoError(t, err)
	})

	t.Run("failed, bank account not found", func(t *testing.T) {
		id := uuid.New()
		userId := uuid.New()
		pgxMock.ExpectQuery("SELECT (.+)").
			WithArgs(id, userId).
			WillReturnError(constant.ErrBankAccountNotFound)

		err := bankaccountSvc.Delete(context.Background(), id, userId)
		assert.Error(t, err)
	})
}
