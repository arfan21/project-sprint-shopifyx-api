package payment

import (
	"context"

	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
)

type Service interface {
	Create(ctx context.Context, req model.PaymentRequest) (err error)
}
