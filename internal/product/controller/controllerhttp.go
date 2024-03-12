package productctrl

import (
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/product"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/exception"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/logger"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ControllerHTTP struct {
	svc product.Service
}

func New(svc product.Service) *ControllerHTTP {
	return &ControllerHTTP{svc: svc}
}

// @Summary Create product
// @Description Create product
// @Tags product
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.ProductRequest true "Payload product create request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/product [post]
func (ctrl ControllerHTTP) Create(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.ProductRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID

	err = ctrl.svc.Create(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "product added successfully",
	})
}

// @Summary Update product
// @Description Update product
// @Tags product
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.ProductRequest true "Payload product update request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/product/{id} [patch]
func (ctrl ControllerHTTP) Update(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}
	var req model.ProductRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID

	idStr := c.Params("id")
	req.ID, err = uuid.Parse(idStr)
	exception.PanicIfNeeded(err)

	err = ctrl.svc.Update(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "product updated successfully",
	})
}

// @Summary Delete product
// @Description Delete product
// @Tags product
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param id path string true "Product ID"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/product/{id} [delete]
func (ctrl ControllerHTTP) Delete(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	exception.PanicIfNeeded(err)

	req := model.ProductDeleteRequest{
		ID:     id,
		UserID: claims.UserID,
	}

	err = ctrl.svc.Delete(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "product deleted successfully",
	})
}

// @Summary Get product list
// @Description Get product list
// @Tags product
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param userOnly query bool false "Get product list by user"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param condition query string false "Condition"
// @Param tags query string false "Tags"
// @Param showEmptyStock query bool false "Show empty stock"
// @Param maxPrice query float64 false "Max price"
// @Param minPrice query float64 false "Min price"
// @Param sortBy query string false "Sort by"
// @Param orderBy query string false "Order by"
// @Param search query string false "Search"
// @Success 200 {object} pkgutil.HTTPResponse{data=[]model.ProductGetResponse}
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/product [get]
func (ctrl ControllerHTTP) GetList(c *fiber.Ctx) error {
	var req model.ProductGetListRequest
	err := c.QueryParser(&req)
	exception.PanicIfNeeded(err)

	res, total, err := ctrl.svc.GetList(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Data: res,
		Meta: pkgutil.MetaResponse{
			Offset: req.Offset,
			Limit:  req.Limit,
			Total:  total,
		},
	})
}

// @Summary Get product detail
// @Description Get product detail
// @Tags product
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param id path string true "Product ID"
// @Success 200 {object} pkgutil.HTTPResponse{data=model.ProductGetResponse}
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/product/{id} [get]
func (ctrl ControllerHTTP) GetDetailByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	exception.PanicIfNeeded(err)

	res, err := ctrl.svc.GetDetailByID(c.UserContext(), id)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Data: res,
	})
}

// @Summary Update stock product
// @Description Update stock product
// @Tags product
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param body body model.ProductUpdateStockRequest true "Payload product update stock request"
// @Success 200 {object} pkgutil.HTTPResponse
// @Failure 400 {object} pkgutil.HTTPResponse{data=[]pkgutil.ErrValidationResponse} "Error validation field"
// @Failure 500 {object} pkgutil.HTTPResponse
// @Router /v1/product/{id}/stock [post]
func (ctrl ControllerHTTP) UpdateStock(c *fiber.Ctx) error {
	claims, ok := c.Locals(constant.JWTClaimsContextKey).(model.JWTClaims)
	if !ok {
		logger.Log(c.UserContext()).Error().Msg("cannot get claims from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid or expired token",
		})
	}

	var req model.ProductUpdateStockRequest
	err := c.BodyParser(&req)
	exception.PanicIfNeeded(err)

	idStr := c.Params("id")
	req.ID, err = uuid.Parse(idStr)
	exception.PanicIfNeeded(err)

	req.UserID = claims.UserID

	err = ctrl.svc.UpdateStock(c.UserContext(), req)
	exception.PanicIfNeeded(err)

	return c.Status(fiber.StatusOK).JSON(pkgutil.HTTPResponse{
		Message: "stock updated successfully",
	})
}
