package server

import (
	bankaccountctrl "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/controller"
	bankaccountrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/repository"
	bankaccountsvc "github.com/arfan21/project-sprint-shopifyx-api/internal/bankaccount/service"
	productctrl "github.com/arfan21/project-sprint-shopifyx-api/internal/product/controller"
	productrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/product/repository"
	productsvc "github.com/arfan21/project-sprint-shopifyx-api/internal/product/service"
	userctrl "github.com/arfan21/project-sprint-shopifyx-api/internal/user/controller"
	userrepo "github.com/arfan21/project-sprint-shopifyx-api/internal/user/repository"
	usersvc "github.com/arfan21/project-sprint-shopifyx-api/internal/user/service"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Routes() {

	api := s.app.Group("")
	api.Get("/health-check", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	userRepo := userrepo.New(s.db)
	userSvc := usersvc.New(userRepo)
	userCtrl := userctrl.New(userSvc)

	productRepo := productrepo.New(s.db)
	productSvc := productsvc.New(productRepo)
	productCtrl := productctrl.New(productSvc)

	bankAccountRepo := bankaccountrepo.New(s.db)
	bankAccountSvc := bankaccountsvc.New(bankAccountRepo)
	bankAccountCtrl := bankaccountctrl.New(bankAccountSvc)

	s.RoutesCustomer(api, userCtrl)
	s.RoutesProduct(api, productCtrl)
	s.RoutesBankAccount(api, bankAccountCtrl)
}

func (s Server) RoutesCustomer(route fiber.Router, ctrl *userctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	usersV1 := v1.Group("/user")
	usersV1.Post("/register", ctrl.Register)
	usersV1.Post("/login", ctrl.Login)
}

func (s Server) RoutesProduct(route fiber.Router, ctrl *productctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	productV1 := v1.Group("/product")
	productV1.Post("", middleware.JWTAuth, ctrl.Create)
	productV1.Patch("/:id", middleware.JWTAuth, ctrl.Update)
	productV1.Delete("/:id", middleware.JWTAuth, ctrl.Delete)
	productV1.Get("", ctrl.GetList)
	productV1.Get("/:id", ctrl.GetDetailByID)
	productV1.Post("/:id/stock", middleware.JWTAuth, ctrl.UpdateStock)
}

func (s Server) RoutesBankAccount(route fiber.Router, ctrl *bankaccountctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	bankAccountV1 := v1.Group("/bank/account", middleware.JWTAuth)
	bankAccountV1.Post("", ctrl.Create)
	bankAccountV1.Patch("/:id", ctrl.Update)
}
