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
