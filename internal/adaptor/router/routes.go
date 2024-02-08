package router

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/tigertony2536/go-login/internal/config"
)

type Router struct {
	app *fiber.App
}

func NewRouter(a *AuthHandler, cfg *config.Config) *Router {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello World")
	})

	app.Post("/login", a.Login(cfg.HttpConfig.JwtSecret))
	app.Get("/refresh", a.Refresh)

	//Middleware 2: Validating JWT
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.HttpConfig.JwtSecret),
	}))

	// app.Get("/logout", )

	//Middleware 2: Check role
	app.Use(a.CheckRole)
	return &Router{app: app}
}

func (r *Router) Serve(cfg *config.HttpConfig) error {
	return r.app.Listen(cfg.SERVER_Port)
}
