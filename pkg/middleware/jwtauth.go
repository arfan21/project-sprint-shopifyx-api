package middleware

import (
	"fmt"
	"strings"

	"github.com/arfan21/project-sprint-shopifyx-api/config"
	"github.com/arfan21/project-sprint-shopifyx-api/internal/model"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/constant"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/logger"
	"github.com/arfan21/project-sprint-shopifyx-api/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JWTAuth(c *fiber.Ctx) error {
	// fetch token
	head := c.Get("Authorization", "")
	if head == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(pkgutil.HTTPResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "missing or malformed jwt",
		})
	}

	token := strings.Split(head, "Bearer ")
	if len(token) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(pkgutil.HTTPResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "missing or malformed jwt",
		})
	}

	// validate token
	t, err := jwt.ParseWithClaims(token[1], &model.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, fmt.Errorf("middleware: invalid token signing algorithm")
		}

		return []byte(config.Get().JWT.Secret), nil
	})
	if err != nil {
		logger.Log(c.UserContext()).Error().Msgf("middleware: failed to parse jwt token: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(pkgutil.HTTPResponse{
			Code:    fiber.StatusUnauthorized,
			Message: "invalid or expired token",
		})
	}

	claims, ok := t.Claims.(*model.JWTClaims)
	if ok && t.Valid && claims != nil {
		idUUID, err := uuid.Parse(claims.Subject)
		if err != nil {
			logger.Log(c.UserContext()).Error().Msgf("middleware: failed to parse user id: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(pkgutil.HTTPResponse{
				Code:    fiber.StatusUnauthorized,
				Message: "invalid or expired token",
			})
		}
		claims.UserID = idUUID
		c.Locals(constant.JWTClaimsContextKey, *claims)
		return c.Next()
	}

	logger.Log(c.UserContext()).Error().Msg("middleware: invalid or expired token")
	return c.Status(fiber.StatusUnauthorized).JSON(pkgutil.HTTPResponse{
		Code:    fiber.StatusUnauthorized,
		Message: "invalid or expired token",
	})
}
