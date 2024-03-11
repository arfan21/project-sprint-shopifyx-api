package usersvc

import (
	"context"
	"testing"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/user"
	userrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/user/repository"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var (
	pgxMock     pgxmock.PgxPoolIface
	redisClient *redis.Client
	redisMock   redismock.ClientMock
	repoPG      user.Repository
	repoRedis   user.RepositoryRedis
	userSvc     *Service

	defaultPassword       = "test123qwe"
	defaultHashedPassword = "$2a$10$BAmWsmtMZvoZcVwjkYpe5uJtuxb/Ii5Il4RHwDTEup9kun6FrZN8."
)

func initDep(t *testing.T) {
	if pgxMock == nil {
		mockPool, err := pgxmock.NewPool()
		assert.NoError(t, err)

		pgxMock = mockPool
	}

	if redisClient == nil || redisMock == nil {
		client, clientMock := redismock.NewClientMock()
		redisClient = client
		redisMock = clientMock
	}

	if repoPG == nil {
		repoPG = userrepo.New(pgxMock)
	}

	if userSvc == nil {
		userSvc = New(repoPG)
	}
}

func TestRegister(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.UserRegisterRequest{
			Name:     "test name",
			Username: "test",
			Password: "test123qwe",
		}

		pgxMock.ExpectExec("INSERT INTO users (.+)").
			WithArgs(pgxmock.AnyArg(), req.Username, req.Name, pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		_, err := userSvc.Register(context.Background(), req)

		assert.NoError(t, err)
	})

	t.Run("failed, Username already registered", func(t *testing.T) {
		req := model.UserRegisterRequest{
			Name:     "test name",
			Username: "test",
			Password: "test123qwe",
		}

		pgxMock.ExpectExec("INSERT INTO users (.+)").
			WithArgs(pgxmock.AnyArg(), req.Username, req.Name, pgxmock.AnyArg()).
			WillReturnError(&pgconn.PgError{Code: "23505"}) // unique violation

		_, err := userSvc.Register(context.Background(), req)

		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrUsernameAlreadyRegistered)
	})

	t.Run("failed, invalid request", func(t *testing.T) {
		req := model.UserRegisterRequest{
			Name:     "test",
			Username: "test",
			Password: "test",
		}

		_, err := userSvc.Register(context.Background(), req)

		assert.Error(t, err)

		var validationErr *constant.ErrValidation
		assert.ErrorAs(t, err, &validationErr)
	})
}

func TestLogin(t *testing.T) {
	initDep(t)

	t.Run("success", func(t *testing.T) {
		req := model.UserLoginRequest{
			Username: "test",
			Password: "test123qwe",
		}

		id := uuid.New()

		pgxMock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(req.Username).
			WillReturnRows(pgxMock.NewRows([]string{"id", "username", "name", "password"}).
				AddRow(id, "test", req.Username, defaultHashedPassword))

		res, err := userSvc.Login(context.Background(), req)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)

		assert.NoError(t, pgxMock.ExpectationsWereMet())
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("failed, username not found", func(t *testing.T) {
		req := model.UserLoginRequest{
			Username: "test",
			Password: "test123qwe",
		}

		pgxMock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(req.Username).
			WillReturnError(pgx.ErrNoRows)

		_, err := userSvc.Login(context.Background(), req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrUsernameOrPasswordInvalid)
	})

	t.Run("failed, invalid password", func(t *testing.T) {
		req := model.UserLoginRequest{
			Username: "test",
			Password: "test123qweasd",
		}
		id := uuid.New()
		pgxMock.ExpectQuery("SELECT (.+) FROM users").
			WithArgs(req.Username).
			WillReturnRows(pgxMock.NewRows([]string{"id", "username", "name", "password"}).
				AddRow(id, "test", req.Username, defaultHashedPassword))

		_, err := userSvc.Login(context.Background(), req)
		assert.Error(t, err)
		assert.ErrorIs(t, err, constant.ErrUsernameOrPasswordInvalid)
	})
}
