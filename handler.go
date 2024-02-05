package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/gofiber/fiber/v2"
)

func Login(secretKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req User
		c.BodyParser(&req)

		if req.Email != user.Email || req.Password != user.Password {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		token, err := CreateToken(req.Email, "users")
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

func GetQuotes(c *fiber.Ctx) error {
	var quotes Quote
	jsonfile, err := os.Open("quotes.json")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "open json file fail",
		})
	}
	data, err := io.ReadAll(jsonfile)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "can not read data from json file",
		})
	}
	err = json.Unmarshal(data, &quotes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unmashal json data fail",
		})
	}
	return c.JSON(fiber.Map{
		"quotes": quotes,
	})
}
