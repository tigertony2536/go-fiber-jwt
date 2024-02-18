package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tigertony2536/go-login/internal/core"
	"github.com/tigertony2536/go-login/internal/core/domain"
)

type AuthHandler struct {
	as *core.AuthServiceImpl
}

func NewAuthHandler(as *core.AuthServiceImpl) *AuthHandler {
	return &AuthHandler{as: as}
}

type Response struct {
	ErrorMessage string
}

func (a *AuthHandler) Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req domain.UserLogin
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(Response{
				ErrorMessage: "username or password wrong",
			})
		}
		err = a.as.Register(&req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(Response{
				ErrorMessage: "This email have been used",
			})
		}
		return c.JSON(req)
	}
}

func (a *AuthHandler) Login(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req domain.UserLogin
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		token, err := a.as.Login(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
		})
	}
}

func (a *AuthHandler) Refresh(c *fiber.Ctx) error {
	rt := c.Get("Authorization")

	return c.JSON(fiber.Map{
		"Refresh_token": rt,
	})
}

const userContextKey = "user"

func (a *AuthHandler) CheckRole(c *fiber.Ctx) error {
	token := c.Locals(userContextKey).(*jwt.Token)
	claim := token.Claims.(jwt.MapClaims)
	if claim["role"] != "users" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}
