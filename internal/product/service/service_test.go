package productsvc

import (
	"context"
	"testing"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount"
	bankaccountrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/repository"
	bankaccountsvc "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/service"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/product"
	productrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/product/repository"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	productSvc      product.Service
	pgxMock         pgxmock.PgxPoolIface
	repoPG          product.Repository
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
		repoPG = productrepo.New(pgxMock)
	}

	if bankAccountRepo == nil {
		bankAccountRepo = bankaccountrepo.New(pgxMock)
	}

	if bankAccountSvc == nil {
		bankAccountSvc = bankaccountsvc.New(bankAccountRepo)
	}

	if productSvc == nil {
		productSvc = New(repoPG, bankAccountSvc)
	}
}

func TestCreate(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.ProductRequest{
			UserID:         uuid.New(),
			Name:           "test name",
			Price:          &decimal.Zero,
			ImageUrl:       "https://test.com/image.jpg",
			Stock:          new(int),
			Condition:      "new",
			Tags:           nil,
			IsPurchaseable: new(bool),
		}

		*req.Stock = 10
		*req.IsPurchaseable = true

		if req.Tags == nil {
			req.Tags = []string{}
		}

		pgxMock.ExpectExec("INSERT INTO products (.+)").
			WithArgs(pgxmock.AnyArg(), req.Name, *req.Price, req.ImageUrl, *req.Stock, entity.ProductCondition(req.Condition), req.Tags, *req.IsPurchaseable, req.UserID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := productSvc.Create(context.Background(), req)

		assert.NoError(t, err)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})

	t.Run("failed, validation error", func(t *testing.T) {
		req := model.ProductRequest{
			Name: "test name",
		}

		err := productSvc.Create(context.Background(), req)

		assert.Error(t, err)
		var validationErr *constant.ErrValidation
		assert.ErrorAs(t, err, &validationErr)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})
}

func TestUpdate(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.ProductRequest{
			ID:             uuid.New(),
			UserID:         uuid.New(),
			Name:           "test name",
			Price:          &decimal.Zero,
			ImageUrl:       "https://test.com/image.jpg",
			Stock:          new(int),
			Condition:      "new",
			Tags:           nil,
			IsPurchaseable: new(bool),
		}

		*req.Stock = 10
		*req.IsPurchaseable = true

		if req.Tags == nil {
			req.Tags = []string{}
		}

		getByIDQueryMock(req.ID, req.UserID)

		pgxMock.ExpectExec("UPDATE products (.+) WHERE id = (.+)").
			WithArgs(req.Name, *req.Price, req.ImageUrl, entity.ProductCondition(req.Condition), req.Tags, *req.IsPurchaseable, req.ID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := productSvc.Update(context.Background(), req)

		assert.NoError(t, err)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})

	t.Run("failed, validation error", func(t *testing.T) {
		req := model.ProductRequest{
			Name: "test name",
		}

		err := productSvc.Update(context.Background(), req)

		assert.Error(t, err)
		var validationErr *constant.ErrValidation
		assert.ErrorAs(t, err, &validationErr)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})

	t.Run("failed, product not found", func(t *testing.T) {
		req := model.ProductRequest{
			ID:             uuid.New(),
			UserID:         uuid.New(),
			Name:           "test name",
			Price:          &decimal.Zero,
			ImageUrl:       "https://test.com/image.jpg",
			Stock:          new(int),
			Condition:      "new",
			Tags:           nil,
			IsPurchaseable: new(bool),
		}

		*req.Stock = 10
		*req.IsPurchaseable = true

		if req.Tags == nil {
			req.Tags = []string{}
		}

		pgxMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+)").
			WithArgs(req.ID).
			WillReturnError(constant.ErrProductNotFound)

		err := productSvc.Update(context.Background(), req)

		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrProductNotFound)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})

	t.Run("failed, user id not match", func(t *testing.T) {
		req := model.ProductRequest{
			ID:             uuid.New(),
			UserID:         uuid.New(),
			Name:           "test name",
			Price:          &decimal.Zero,
			ImageUrl:       "https://test.com/image.jpg",
			Stock:          new(int),
			Condition:      "new",
			Tags:           nil,
			IsPurchaseable: new(bool),
		}

		*req.Stock = 10
		*req.IsPurchaseable = true

		if req.Tags == nil {
			req.Tags = []string{}
		}

		getByIDQueryMock(req.ID, uuid.New())

		err := productSvc.Update(context.Background(), req)

		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrAccessForbidden)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})
}

func TestDelete(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		userId := uuid.New()

		req := model.ProductDeleteRequest{
			ID:     id,
			UserID: userId,
		}

		getByIDQueryMock(id, userId)

		pgxMock.ExpectExec("DELETE FROM products WHERE id = (.+)").
			WithArgs(req.ID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err := productSvc.Delete(context.Background(), req)

		assert.NoError(t, err)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})

}

func getByIDQueryMock(id, userId uuid.UUID) {
	pgxMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+)").
		WithArgs(id).
		WillReturnRows(pgxMock.NewRows([]string{"id", "name", "price", "imageUrl", "stock", "condition", "tags", "isPurchaseable", "userId"}).
			AddRow(id, "test name", decimal.Zero, "https://test.com/image.jpg", 10, entity.ProductCondition("new"), nil, true, userId))
}

func TestGetList(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.ProductGetListRequest{
			Offset: 0,
			Limit:  10,
		}

		productId := uuid.New()

		pgxMock.ExpectQuery("SELECT (.+) FROM products (.+)").
			WithArgs(0, req.Limit, req.Offset).
			WillReturnRows(pgxMock.NewRows([]string{"id", "name", "price", "imageUrl", "stock", "condition", "tags", "isPurchaseable"}).
				AddRow(productId, "test name", decimal.Zero, "https://test.com/image.jpg", 10, entity.ProductCondition("new"), nil, true))

		pgxMock.ExpectQuery("SELECT COUNT(.+) FROM products (.+)").
			WithArgs(0).
			WillReturnRows(pgxMock.NewRows([]string{"count"}).AddRow(1))

		pgxMock.ExpectQuery("SELECT (.+) FROM payments (.+)").
			WithArgs([]uuid.UUID{productId}).
			WillReturnRows(pgxMock.NewRows([]string{"total", "productId"}).AddRow(1, productId))

		res, total, err := productSvc.GetList(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, 1, total)
		assert.Equal(t, 1, len(res))
		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})
}

func TestGetDetail(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		id := uuid.New()

		pgxMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+)").
			WithArgs(id).
			WillReturnRows(pgxMock.NewRows([]string{"id", "name", "price", "imageUrl", "stock", "condition", "tags", "isPurchaseable", "userId"}).
				AddRow(id, "test name", decimal.Zero, "https://test.com/image.jpg", 10, entity.ProductCondition("new"), nil, true, uuid.New()))

		res, err := productSvc.GetDetailByID(context.Background(), id)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})

	t.Run("failed, product not found", func(t *testing.T) {
		id := uuid.New()

		pgxMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+)").
			WithArgs(id).
			WillReturnError(constant.ErrProductNotFound)

		_, err := productSvc.GetDetailByID(context.Background(), id)

		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrProductNotFound)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})
}

func TestUpdateStock(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.ProductUpdateStockRequest{
			ID:     uuid.New(),
			UserID: uuid.New(),
			Stock:  10,
		}

		getByIDQueryMock(req.ID, req.UserID)

		pgxMock.ExpectBegin()
		pgxMock.ExpectQuery("SELECT (.+) FROM products WHERE id = (.+) FOR UPDATE").
			WithArgs(req.ID).
			WillReturnRows(pgxMock.NewRows([]string{"stock"}).
				AddRow(10))

		pgxMock.ExpectExec("UPDATE products (.+) WHERE id = (.+)").
			WithArgs(req.Stock, req.ID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))
		pgxMock.ExpectCommit()

		err := productSvc.UpdateStock(context.Background(), req)

		assert.NoError(t, err)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})

}

func TestReduceStock(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		qty := 10
		id := uuid.New()

		pgxMock.ExpectExec("UPDATE products (.+) WHERE id = (.+) AND (.+)").
			WithArgs(qty, id).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := productSvc.ReduceStock(context.Background(), id, qty)

		assert.NoError(t, err)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})

	t.Run("failed, stock not enough", func(t *testing.T) {
		qty := 10
		id := uuid.New()

		pgxMock.ExpectExec("UPDATE products (.+) WHERE id = (.+) AND (.+)").
			WithArgs(qty, id).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := productSvc.ReduceStock(context.Background(), id, qty)

		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrInsufficientStock)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})
}

func TestPayment(t *testing.T) {
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

		err := productSvc.Payment(context.Background(), req)
		assert.NoError(t, err)
		assert.NoError(t, pgxMock.ExpectationsWereMet())
	})

	t.Run("failed, validation error", func(t *testing.T) {
		req := model.PaymentRequest{
			ProductID: uuid.New(),
			UserID:    uuid.New(),
		}

		err := productSvc.Payment(context.Background(), req)
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

		err := productSvc.Payment(context.Background(), req)
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

		err := productSvc.Payment(context.Background(), req)
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

		err := productSvc.Payment(context.Background(), req)
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

		err := productSvc.Payment(context.Background(), req)
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

		err := productSvc.Payment(context.Background(), req)
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

		err := productSvc.Payment(context.Background(), req)
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
