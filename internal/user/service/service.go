package usersvc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/arfan21/project-sprint-shopifyx-api/config"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/entity"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/user"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/validation"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo user.Repository
}

func New(repo user.Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) Register(ctx context.Context, req model.UserRegisterRequest) (res model.UserLoginResponse, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to validate request: %w", err)
		return
	}

	cost := bcrypt.DefaultCost
	if config.Get().Bcrypt.Salt > 0 {
		cost = config.Get().Bcrypt.Salt
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), cost)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to hash password: %w", err)
		return
	}

	data := entity.User{
		Name:     req.Name,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	err = s.repo.Create(ctx, data)
	if err != nil {
		err = fmt.Errorf("user.service.Register: failed to register user: %w", err)
		return
	}

	accessTokenExpire := time.Duration(config.Get().JWT.ExpireIn) * time.Second

	accessToken, err := s.CreateJWTWithExpiry(
		data.ID.String(),
		data.Name,
		config.Get().JWT.Secret,
		accessTokenExpire,
	)

	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to create access token: %w", err)
		return
	}
	res = model.UserLoginResponse{
		Username:    data.Username,
		Name:        data.Name,
		AccessToken: accessToken,
	}
	return
}

func (s Service) Login(ctx context.Context, req model.UserLoginRequest) (res model.UserLoginResponse, err error) {
	err = validation.Validate(req)
	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to validate request: %w", err)
		return
	}

	data, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to get user by email: %w", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = constant.ErrUsernameOrPasswordInvalid
		}
		err = fmt.Errorf("user.service.Login: failed to compare password: %w", err)
		return
	}

	accessTokenExpire := time.Duration(config.Get().JWT.ExpireIn) * time.Second

	accessToken, err := s.CreateJWTWithExpiry(
		data.ID.String(),
		data.Name,
		config.Get().JWT.Secret,
		accessTokenExpire,
	)

	if err != nil {
		err = fmt.Errorf("user.service.Login: failed to create access token: %w", err)
		return
	}

	res = model.UserLoginResponse{
		Username:    data.Username,
		Name:        data.Name,
		AccessToken: accessToken,
	}

	return
}

func (s Service) CreateJWTWithExpiry(id, name, secret string, expiry time.Duration) (token string, err error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, model.JWTClaims{
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.Get().Service.Name,
			Subject:   id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	token, err = jwtToken.SignedString([]byte(secret))
	if err != nil {
		err = fmt.Errorf("usecase: failed to create jwt token: %w", err)
		return
	}

	return
}
