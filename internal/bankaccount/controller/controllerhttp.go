package bankaccountctrl

import (
	"github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/exception"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/logger"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
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

// @Summary Update Bank Account
// @Description Update Bank Account
// @Tags Bank Account
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.BankAccountRequest true "Payload bank account update request"
// @Param id path string true "Bank Account ID"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/bank/account/{id} [patch]
func (ctrl ControllerHTTP) Update(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.BankAccountRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	exception.PanicIfNeeded(err)
	req.BankAccountID = id

	err = ctrl.svc.Update(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "bank account updated successfully",
	})
}

// @Summary Delete Bank Account
// @Description Delete Bank Account
// @Tags Bank Account
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param id path string true "Bank Account ID"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/bank/account/{id} [delete]
func (ctrl ControllerHTTP) Delete(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	exception.PanicIfNeeded(err)

	err = ctrl.svc.Delete(c.UserContext(), id, claims.UserID)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "bank account deleted successfully",
	})
}

// @Summary Get Bank Accounts
// @Description Get Bank Accounts
// @Tags Bank Account
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Success 200 {object} pkgutil.HTTPResponse{data=[]model.BankAccountResponse}
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/bank/account [get]
func (ctrl ControllerHTTP) GetList(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	res, err := ctrl.svc.GetListByUserID(c.UserContext(), claims.UserID)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "success",
		Data:    res,
	})
}
