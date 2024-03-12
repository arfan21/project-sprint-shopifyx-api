package bankaccountctrl

import (
	"github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/exception"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/logger"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type ControllerHTTP struct {
	svc bankaccount.Service
}

func New(svc bankaccount.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Create Bank Account
// @Description Create Bank Account
// @Tags Bank Account
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.BankAccountRequest true "Payload bank account create request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1//v1/bank/account [post]
func (ctrl ControllerHTTP) Create(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.BankAccountRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID

	err = ctrl.svc.Create(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "bank account added successfully",
	})
}
