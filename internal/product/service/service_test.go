package productsvc

import (
	"context"
	"testing"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/product"
	productrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/product/repository"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var (
	productSvc product.Service
	pgxMock    pgxmock.PgxPoolIface
	repoPG     product.Repository
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

	if productSvc == nil {
		productSvc = New(repoPG)
	}
}

func TestCreate(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.ProductCreateRequest{
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
			WithArgs(req.Name, *req.Price, req.ImageUrl, *req.Stock, entity.ProductCondition(req.Condition), req.Tags, *req.IsPurchaseable, req.UserID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := productSvc.Create(context.Background(), req)

		assert.NoError(t, err)
	})

	t.Run("failed, validation error", func(t *testing.T) {
		req := model.ProductCreateRequest{
			Name: "test name",
		}

		err := productSvc.Create(context.Background(), req)

		assert.Error(t, err)
		var validationErr *constant.ErrValidation
		assert.ErrorAs(t, err, &validationErr)
	})
}
