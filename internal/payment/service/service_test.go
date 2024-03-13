package paymentsvc

import (
	"context"
	"testing"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount"
	bankaccountrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/repository"
	bankaccountsvc "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/service"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/payment"
	paymentrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/payment/repository"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/product"
	productrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/product/repository"
	productsvc "github.com/arfan21/project-sprint-shopifyx-api/internal/product/service"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	paymentSvc      payment.Service
	pgxMock         pgxmock.PgxPoolIface
	repoPG          payment.Repository
	productRepo     product.Repository
	productSvc      product.Service
	bankAccountRepo bankaccount.Repository
	bankAccountSvc  bankaccount.Service
)

func initDep(t *testing.T) {
	if pgxMock == nil {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)

		pgxMock = mockPool
	}

	if repoPG == nil {
		repoPG = paymentrepo.New(pgxMock)
	}

	if productRepo == nil {
		productRepo = productrepo.New(pgxMock)

	}

	if productSvc == nil {
		productSvc = productsvc.New(productRepo)
	}

	if bankAccountRepo == nil {
		bankAccountRepo = bankaccountrepo.New(pgxMock)
	}

	if bankAccountSvc == nil {
		bankAccountSvc = bankaccountsvc.New(bankAccountRepo)
	}

	if paymentSvc == nil {
		paymentSvc = New(repoPG, bankAccountSvc, productSvc)
	}
}

func TestCreate(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.PaymentRequest{
			ProductID:            uuid.New(),
			UserID:               uuid.New(),
			BankAccountID:        uuid.New(),
			Quantity:             1,
			PaymentProofImageURL: "https://test.com/image.jpg",
		}

		getBankAccountByIDQueryMock(req.BankAccountID, req.UserID)
		getProductByIDQueryMock(req.ProductID, uuid.New())

		pgxMock.ExpectBegin()
		pgxMock.ExpectExec("INSERT INTO payments (.+)").
			WithArgs(pgxmock.AnyArg(), req.UserID, req.ProductID, req.BankAccountID, req.PaymentProofImageURL, req.Quantity, decimal.Zero).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		pgxMock.ExpectExec("UPDATE products (.+) WHERE id = (.+) AND (.+)").
			WithArgs(req.Quantity, req.ProductID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		pgxMock.ExpectCommit()

		err := paymentSvc.Create(context.Background(), req)
		assert.NoError(t, err)
		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})

	t.Run("failed, validation error", func(t *testing.T) {
		req := model.PaymentRequest{
			ProductID: uuid.New(),
			UserID:    uuid.New(),
		}

		err := paymentSvc.Create(context.Background(), req)
		assert.Error(t, err)
	})

	t.Run("failed, get bank account error", func(t *testing.T) {
		req := model.PaymentRequest{
			ProductID:            uuid.New(),
			UserID:               uuid.New(),
			BankAccountID:        uuid.New(),
			Quantity:             1,
			PaymentProofImageURL: "https://test.com/image.jpg",
		}

		pgxMock.ExpectQuery("SELECT (.+)").
			WithArgs(req.BankAccountID, req.UserID).
			WillReturnError(pgx.ErrNoRows)

		err := paymentSvc.Create(context.Background(), req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrBankAccountNotFound)
	})

	t.Run("failed, get product error", func(t *testing.T) {
		req := model.PaymentRequest{
			ProductID:            uuid.New(),
			UserID:               uuid.New(),
			BankAccountID:        uuid.New(),
			Quantity:             1,
			PaymentProofImageURL: "https://test.com/image.jpg",
		}

		getBankAccountByIDQueryMock(req.BankAccountID, req.UserID)
		pgxMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+)").
			WithArgs(req.ProductID).
			WillReturnError(pgx.ErrNoRows)

		err := paymentSvc.Create(context.Background(), req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrProductNotFound)
	})

	t.Run("failed, owned product", func(t *testing.T) {
		req := model.PaymentRequest{
			ProductID:            uuid.New(),
			UserID:               uuid.New(),
			BankAccountID:        uuid.New(),
			Quantity:             1,
			PaymentProofImageURL: "https://test.com/image.jpg",
		}

		getBankAccountByIDQueryMock(req.BankAccountID, req.UserID)
		getProductByIDQueryMock(req.ProductID, req.UserID)

		err := paymentSvc.Create(context.Background(), req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrCannotButOwnProduct)
	})

	t.Run("failed, insufficient stock", func(t *testing.T) {
		req := model.PaymentRequest{
			ProductID:            uuid.New(),
			UserID:               uuid.New(),
			BankAccountID:        uuid.New(),
			Quantity:             2,
			PaymentProofImageURL: "https://test.com/image.jpg",
		}

		getBankAccountByIDQueryMock(req.BankAccountID, req.UserID)
		pgxMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+)").
			WithArgs(req.ProductID).
			WillReturnRows(pgxMock.NewRows([]string{"id", "name", "price", "imageUrl", "stock", "condition", "tags", "isPurchaseable", "userId"}).
				AddRow(req.ProductID, "test name", decimal.Zero, "https://test.com/image.jpg", 0, entity.ProductCondition("new"), nil, true, uuid.New()))

		err := paymentSvc.Create(context.Background(), req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrInsufficientStock)
	})

	t.Run("failed, product not purchaseable", func(t *testing.T) {
		req := model.PaymentRequest{
			ProductID:            uuid.New(),
			UserID:               uuid.New(),
			BankAccountID:        uuid.New(),
			Quantity:             1,
			PaymentProofImageURL: "https://test.com/image.jpg",
		}

		getBankAccountByIDQueryMock(req.BankAccountID, req.UserID)
		pgxMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+)").
			WithArgs(req.ProductID).
			WillReturnRows(pgxMock.NewRows([]string{"id", "name", "price", "imageUrl", "stock", "condition", "tags", "isPurchaseable", "userId"}).
				AddRow(req.ProductID, "test name", decimal.Zero, "https://test.com/image.jpg", 10, entity.ProductCondition("new"), nil, false, uuid.New()))

		err := paymentSvc.Create(context.Background(), req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrProductNotPurchaseable)
	})

	t.Run("failed, insufficient stock on reduce stock", func(t *testing.T) {
		req := model.PaymentRequest{
			ProductID:            uuid.New(),
			UserID:               uuid.New(),
			BankAccountID:        uuid.New(),
			Quantity:             2,
			PaymentProofImageURL: "https://test.com/image.jpg",
		}

		getBankAccountByIDQueryMock(req.BankAccountID, req.UserID)
		getProductByIDQueryMock(req.ProductID, uuid.New())

		pgxMock.ExpectBegin()
		pgxMock.ExpectExec("INSERT INTO payments (.+)").
			WithArgs(pgxmock.AnyArg(), req.UserID, req.ProductID, req.BankAccountID, req.PaymentProofImageURL, req.Quantity, decimal.Zero).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		pgxMock.ExpectExec("UPDATE products (.+) WHERE id = (.+) AND (.+)").
			WithArgs(req.Quantity, req.ProductID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))
		pgxMock.ExpectRollback()

		err := paymentSvc.Create(context.Background(), req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrInsufficientStock)
	})
}

func getProductByIDQueryMock(id, userId uuid.UUID) {
	pgxMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+)").
		WithArgs(id).
		WillReturnRows(pgxMock.NewRows([]string{"id", "name", "price", "imageUrl", "stock", "condition", "tags", "isPurchaseable", "userId"}).
			AddRow(id, "test name", decimal.Zero, "https://test.com/image.jpg", 10, entity.ProductCondition("new"), nil, true, userId))
}

func getBankAccountByIDQueryMock(id, userId uuid.UUID) {
	pgxMock.ExpectQuery("SELECT (.+)").
		WithArgs(id, userId).
		WillReturnRows(pgxmock.NewRows([]string{"id", "accountNumber", "accountHolder", "bankName", "userId"}).
			AddRow(id, "1234567890", "Test", "BCA", userId))
}
