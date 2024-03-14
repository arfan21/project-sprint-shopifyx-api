package productctrl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/arfan21/project-sprint-shopifyx-api/config"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount"
	bankaccountrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/repository"
	bankaccountsvc "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/service"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/product"
	productrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/product/repository"
	productsvc "github.com/arfan21/project-sprint-shopifyx-api/internal/product/service"
	usersvc "github.com/arfan21/project-sprint-shopifyx-api/internal/user/service"
	"github.com/arfan21/project-sprint-shopifyx-api/migration"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/exception"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/middleware"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var db *pgxpool.Pool
var dockerPool *dockertest.Pool
var dockerResource *dockertest.Resource

func initDocker(t *testing.T) (*dockertest.Pool, *dockertest.Resource) {
	ctx := context.Background()
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Errorf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		t.Errorf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16.0-alpine3.18",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=postgres-test",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	},
	)
	if err != nil {
		t.Errorf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		connString := fmt.Sprintf("postgres://postgres:postgres@localhost:%s/postgres-test?sslmode=disable", resource.GetPort("5432/tcp"))
		db, err = pgxpool.New(ctx, connString)
		if err != nil {
			return err
		}
		err = db.Ping(ctx)
		return err
	}); err != nil {
		t.Errorf("Could not connect to database: %s", err)
		db = nil
	}

	resource.Expire(30)
	fmt.Println("Database connected")

	dbSql := stdlib.OpenDBFromPool(db)

	fmt.Println("migrate main")
	mig, err := migration.New(dbSql)
	if err != nil {
		t.Errorf("Could not get migrations: %s", err)
	}

	err = mig.Up()
	assert.NoError(t, err)

	return pool, resource
}

var (
	productRepo     product.Repository
	productSvc      product.Service
	bankAccountRepo bankaccount.Repository
	bankAccountSvc  bankaccount.Service
	productCtrl     *ControllerHTTP

	defaultPassword = "test"
)

func initDep(t *testing.T) {
	if dockerPool == nil {
		initDocker(t)
	}

	if bankAccountRepo == nil {
		bankAccountRepo = bankaccountrepo.New(db)
	}

	if bankAccountSvc == nil {
		bankAccountSvc = bankaccountsvc.New(bankAccountRepo)
	}

	if productRepo == nil {
		productRepo = productrepo.New(db)
	}

	if productSvc == nil {
		productSvc = productsvc.New(productRepo, bankAccountSvc)
	}

	if productCtrl == nil {
		productCtrl = New(productSvc)
	}
}

func initUser(t *testing.T) (id uuid.UUID) {
	id = uuid.New()

	query := `
		INSERT INTO users (id, username, name, password)
		VALUES ($1, $2, $3, $4)
	`

	usernameRand := fmt.Sprintf("test-%s", id.String())
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	_, err = db.Exec(context.Background(), query, id, usernameRand, "test", string(hashedPassword))
	assert.NoError(t, err)

	return id
}

func initLogin(t *testing.T, id uuid.UUID) (token string) {
	secret := config.Get().JWT.Secret
	token, err := usersvc.Service{}.CreateJWTWithExpiry(id.String(), "test", secret, 100*time.Second)
	assert.NoError(t, err)

	return token
}

func initBankAccount(t *testing.T, userID uuid.UUID) (id uuid.UUID) {
	id = uuid.New()

	query := `
		INSERT INTO bank_accounts (id, accountNumber, accountHolder, bankName, userId)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := db.Exec(context.Background(), query, id, "111", "test", "test", userID)
	assert.NoError(t, err)

	return id
}

func initProduct(t *testing.T, userID uuid.UUID, stock int) (id uuid.UUID) {
	id = uuid.New()

	query := `
		INSERT INTO products (id, name, price, imageUrl, stock, condition, tags, isPurchaseable, userId)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := db.Exec(context.Background(), query, id, "test", 1000, "test", stock, "new", []string{"test"}, true, userID)
	assert.NoError(t, err)

	return id

}

func getStockProduct(t *testing.T, id uuid.UUID) (stock int) {
	query := `
		SELECT stock
		FROM products
		WHERE id = $1
	`

	err := db.QueryRow(context.Background(), query, id).Scan(&stock)
	assert.NoError(t, err)

	return stock
}

func TestCreateRoute(t *testing.T) {
	initDep(t)
	route := "/v1/product/:id/buy"

	app := fiber.New(fiber.Config{
		ErrorHandler: exception.FiberErrorHandler,
	})

	app.Use(recover.New())
	app.Post(route, middleware.JWTAuth, productCtrl.Payment)

	t.Run("success", func(t *testing.T) {
		sellerId := initUser(t)
		buyerId := initUser(t)
		buyerToken := initLogin(t, buyerId)

		bankAccountId := initBankAccount(t, buyerId)
		productId := initProduct(t, sellerId, 10)

		req := model.PaymentRequest{
			ProductID:            productId,
			UserID:               buyerId,
			BankAccountID:        bankAccountId,
			Quantity:             1,
			PaymentProofImageURL: "https://test.com/image.jpg",
		}
		reqTest := createRequestTest(t, req, route, buyerToken)

		resp, err := app.Test(reqTest, 1000)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		assert.Equal(t, 9, getStockProduct(t, productId))
	})

	t.Run("test concurrent buy", func(t *testing.T) {
		sellerId := initUser(t)
		buyerId := initUser(t)
		buyerToken := initLogin(t, buyerId)
		totalConcurrent := 10

		bankAccountId := initBankAccount(t, buyerId)
		productId := initProduct(t, sellerId, totalConcurrent+1)

		wg := &sync.WaitGroup{}
		arrError := make([]error, totalConcurrent)

		for i := 0; i < totalConcurrent; i++ {
			wg.Add(1)

			go func(wg *sync.WaitGroup, i int) {
				defer wg.Done()
				req := model.PaymentRequest{
					ProductID:            productId,
					UserID:               buyerId,
					BankAccountID:        bankAccountId,
					Quantity:             2,
					PaymentProofImageURL: "https://test.com/image.jpg",
				}
				reqTest := createRequestTest(t, req, route, buyerToken)

				resp, err := app.Test(reqTest, 1000)

				assert.NoError(t, err)

				bodyByte, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var respBody pkgutil.HTTPResponse
				err = json.Unmarshal(bodyByte, &respBody)
				assert.NoError(t, err)

				if resp.StatusCode != 200 {
					arrError[i] = fmt.Errorf(respBody.Message)
				}
			}(wg, i)
		}

		wg.Wait()
		noErrorCounter := 0
		for i, v := range arrError {
			if v != nil {
				fmt.Printf("concurent %d: %s\n", i, v.Error())
				assert.Equal(t, constant.ErrInsufficientStock.Error(), v.Error())
			} else {
				fmt.Printf("concurent %d: success\n", i)
				noErrorCounter++
			}
		}

		assert.Equal(t, totalConcurrent/2, noErrorCounter)

		assert.Equal(t, 1, getStockProduct(t, productId))
	})
}

func createRequestTest(t *testing.T, req model.PaymentRequest, route string, token string) *http.Request {
	reqJson, err := json.Marshal(req)
	assert.NoError(t, err)
	body := new(bytes.Buffer)
	body.Write(reqJson)

	route = fmt.Sprintf("/v1/product/%s/buy", req.ProductID.String())
	reqTest := httptest.NewRequest("POST", route, body)
	reqTest.Header.Set("Content-Type", "application/json")
	reqTest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return reqTest
}
