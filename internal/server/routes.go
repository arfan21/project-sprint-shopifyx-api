package server

import (
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

	s.RoutesCustomer(api, userCtrl)
	s.RoutesProduct(api, productCtrl)
}

func (s Server) RoutesCustomer(route fiber.Router, ctrl *userctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	usersV1 := v1.Group("/user")
	usersV1.Post("/register", ctrl.Register)
	usersV1.Post("/login", ctrl.Login)
}

func (s Server) RoutesProduct(route fiber.Router, ctrl *productctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	productV1 := v1.Group("/product", middleware.JWTAuth)
	productV1.Post("", ctrl.Create)
}
