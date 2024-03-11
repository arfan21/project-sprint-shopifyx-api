package server

import (
	userctrl "github.com/arfan21/shopifyx-api/internal/user/controller"
	userrepo "github.com/arfan21/shopifyx-api/internal/user/repository"
	usersvc "github.com/arfan21/shopifyx-api/internal/user/service"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Routes() {

	api := s.app.Group("/api")
	api.Get("/health-check", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	userRepo := userrepo.New(s.db)
	userSvc := usersvc.New(userRepo)
	userCtrl := userctrl.New(userSvc)

	s.RoutesCustomer(api, userCtrl)

}

func (s Server) RoutesCustomer(route fiber.Router, ctrl *userctrl.ControllerHTTP) {
	v1 := route.Group("/v1")
	usersV1 := v1.Group("/users")
	usersV1.Post("/register", ctrl.Register)
	usersV1.Post("/login", ctrl.Login)
}
