package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tigertony2536/go-login/internal/domain"
	"github.com/tigertony2536/go-login/internal/port"
)

type AuthHandler struct {
	repo port.UserRepository
}

func NewAuthHandler(repo port.UserRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func (a *AuthHandler) Login(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req domain.UserLogin
		c.BodyParser(&req)
		u, err := a.repo.GetUserByEmail(req.Email)
		if err != nil {
			c.JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "user not found",
			})
		}
		if req.Password != u.Password {
			return c.JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "wrong password",
			})
		}
		token, err := domain.CreateToken(req.Email, "users")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "create token fail",
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
