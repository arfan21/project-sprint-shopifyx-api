package user

import (
	"context"
	"time"

	"github.com/arfan21/shopifyx-api/internal/entity"
	userrepo "github.com/arfan21/shopifyx-api/internal/user/repository"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	Begin(ctx context.Context) (tx pgx.Tx, err error)
	WithTx(tx pgx.Tx) *userrepo.Repository

	Create(ctx context.Context, data entity.User) (err error)
	GetByUsername(ctx context.Context, username string) (data entity.User, err error)
}

type RepositoryRedis interface {
	SetRefreshToken(ctx context.Context, token string, expireIn time.Duration, payload entity.UserRefreshToken) (err error)
	IsRefreshTokenExist(ctx context.Context, token string) (payload entity.UserRefreshToken, err error)
	DeleteRefreshToken(ctx context.Context, token string) (err error)
}
