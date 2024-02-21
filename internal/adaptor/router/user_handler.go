package router

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tigertony2536/go-login/internal/adaptor/database"
)

type UserHandler struct {
	Ur *database.UserRepositoryImpl
}

func NewUserHandler(ur *database.UserRepositoryImpl) *UserHandler {
	return &UserHandler{Ur: ur}
}

func (u *UserHandler) GetUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := u.Ur.GetUsers()
		if err != nil {
			log.Fatal("Get users faild. Can not Get user from DB", err)
		}
		return c.Status(fiber.StatusOK).JSON(users)
	}
}
