package router

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/tigertony2536/go-login/internal/config"
)

type Router struct {
	app *fiber.App
}

func NewRouter(a *AuthHandler, u *UserHandler, cfg *config.Config) *Router {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON("Hello World")
	})
	app.Post("/register", a.Register("user"))
	app.Post("/admin", a.Register("admin"))
	app.Post("/login", a.Login(cfg.HttpConfig.JwtSecret))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(cfg.HttpConfig.JwtSecret),
	}))
	app.Get("/refresh", a.Refresh(cfg.HttpConfig.JwtSecret))
	app.Post("/home", func(c *fiber.Ctx) error {
		return c.JSON("Welcome to Homepage")
	})
	app.Delete("/delete", a.DeleteAccount())
	app.Use(a.Protected(cfg.HttpConfig.JwtSecret))
	app.Get("/users", u.GetUsers())
	//Middleware 2: Check role

	return &Router{app: app}
}

func (r *Router) Serve(cfg *config.HttpConfig) error {
	return r.app.Listen(cfg.SERVER_Port)
}
