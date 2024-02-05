package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
)

var (
	secretKey = os.Getenv("JWT_SECRET")
	user      = User{
		Email:    "tigertony2536@gmail.com",
		Password: "123456",
	}
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(user)
	})

	app.Post("/login", Login(secretKey))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(secretKey),
	}))
	// app.Use(extractUserFromJWT)
	app.Use(checkMiddleware)
	app.Get("/quotes", GetQuotes)
	app.Listen(":8000")

}
