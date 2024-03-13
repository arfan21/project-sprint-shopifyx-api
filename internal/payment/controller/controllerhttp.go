package paymentctrl

import (
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/payment"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/exception"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/logger"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type ControllerHTTP struct {
	svc payment.Service
}

func New(svc payment.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Create Payment
// @Description Create Payment
// @Tags Payment
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.PaymentRequest true "Payload payment create request"
// @Param id path string true "Product ID"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/product/{id}/buy [post]
func (ctrl ControllerHTTP) Create(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.PaymentRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID

	err = ctrl.svc.Create(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "payment added successfully",
	})
}
