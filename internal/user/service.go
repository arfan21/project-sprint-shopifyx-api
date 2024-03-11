package user

import (
	"context"

	"github.com/arfan21/shopifyx-api/internal/model"
)

type Service interface {
	Register(ctx context.Context, req model.UserRegisterRequest) (err error)
	Login(ctx context.Context, req model.UserLoginRequest) (res model.UserLoginResponse, err error)
}
