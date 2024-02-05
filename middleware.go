package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const userContextKey = "user"

func extractUserFromJWT(c *fiber.Ctx) error {
	user := &Token{}

	// Extract the token from the Fiber context (inserted by the JWT middleware)
	token := c.Locals(userContextKey).(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	fmt.Println(claims)

	user.Email = claims["email"].(string)
	user.Role = claims["role"].(string)

	// Store the user data in the Fiber context
	c.Locals(userContextKey, user)

	return c.Next()
}

func checkMiddleware(c *fiber.Ctx) error {
	token := c.Locals(userContextKey).(*jwt.Token)
	claim := token.Claims.(jwt.MapClaims)
	if claim["role"] != "users" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	return c.Next()
}
