package router

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tigertony2536/go-login/internal/adaptor/database"
	"github.com/tigertony2536/go-login/internal/core"
	"github.com/tigertony2536/go-login/internal/core/domain"
)

type AuthHandler struct {
	As *core.AuthServiceImpl
	Ar *database.AuthRepositoryImpl
}

func NewAuthHandler(as *core.AuthServiceImpl, ar *database.AuthRepositoryImpl) *AuthHandler {
	return &AuthHandler{As: as, Ar: ar}
}

type Response struct {
	ErrorMessage string
}

func (a *AuthHandler) Register(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req domain.UserLogin
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(Response{
				ErrorMessage: "username or password wrong",
			})
		}
		err = a.As.Register(&req, role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(Response{
				ErrorMessage: "This email have been used",
			})
		}
		return c.JSON(req)
	}
}

func (a *AuthHandler) Login(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//Log in with co
		var req domain.UserLogin
		err := c.BodyParser(&req)
		if err != nil {
			return c.JSON(fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "Login failed" + err.Error(),
			})
		}
		session, err := a.Ar.GetSessionByUserID(req.ID)
		if err == nil || session != nil {
			expired, err := a.As.IsExpired(session.CreatedAt)
			if err != nil {
				return c.JSON(fiber.Error{
					Code:    fiber.StatusBadRequest,
					Message: "Login failed" + err.Error(),
				})
			}
			if !expired {
				return c.SendString("You have Logged In")
			}
		}
		user, err := a.As.Repo.GetUserByEmail(req.Email)
		if err != nil {
			return c.JSON(fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "Login failed" + err.Error(),
			})
		}
		newToken, err := a.As.Login(user, secret)
		if err != nil {
			return c.JSON(fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "Login failed" + err.Error(),
			})
		}
		s := domain.Session{UserID: user.ID, RefreshToken: newToken.RefreshToken, CreatedAt: time.Now()}
		err = a.Ar.CreateSession(&s)
		if err != nil {
			return c.JSON(fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "Login failed" + err.Error(),
			})
		}

		ses, err := a.Ar.GetSessionByUserID(s.UserID)
		if err != nil {
			return c.JSON(fiber.Error{
				Code:    fiber.StatusBadRequest,
				Message: "Login failed" + err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"accessToken":  newToken.AccessToken,
			"refreshToken": newToken.RefreshToken,
			"sessionID":    ses.SessionID,
			"userID":       ses.UserID,
		})
	}
}

func (a *AuthHandler) Refresh(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claim := token.Claims.(jwt.MapClaims)
		if claim.Valid() != nil {
			return c.JSON(fiber.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "Refresh token expired",
			})
		}
		id := claim["UserID"].(float64)
		uid := uint(id)
		newToken, err := a.As.Refresh(uid, token, secret)
		if err != nil {
			return c.JSON(fiber.Error{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"access_token":  newToken.AccessToken,
			"refresh_token": newToken.RefreshToken,
		})
	}
}

func (a *AuthHandler) Protected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claim := token.Claims.(jwt.MapClaims)
		if claim.Valid() != nil {
			return c.Redirect("/", fiber.StatusUnauthorized)
		}
		role := claim["Role"].(string)
		fmt.Println(role)
		if role != "admin" {
			return c.JSON(fiber.Error{
				Code:    fiber.StatusUnauthorized,
				Message: "User can not access this route.",
			})
		}
		return c.Next()
	}
}

func (a *AuthHandler) DeleteAccount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Locals("user").(*jwt.Token)
		claim := token.Claims.(jwt.MapClaims)
		email := claim["Email"].(string)
		err := a.As.Repo.DeleteUser(email)
		if err != nil {
			return c.JSON(fiber.Error{
				Code:    fiber.StatusInternalServerError,
				Message: "Delete user failed" + err.Error(),
			})
		}
		u, err := a.As.Repo.GetUserByEmail(email)
		if err != nil {
			return c.JSON(fiber.Error{
				Code:    fiber.StatusInternalServerError,
				Message: "Delete user failed" + err.Error(),
			})
		}
		if u != nil {
			return c.SendString("Delete user failed")
		}
		return c.SendString("Delete user successfully")
	}
}
