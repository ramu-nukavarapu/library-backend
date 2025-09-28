package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get the token from either cookies or request headers
		cookieToken := c.Cookies("jwt")
		var token string

		if cookieToken != "" {
			log.Warn("token from cookies, using it...")
			token = cookieToken
		} else {
			log.Warn("empty token from cookies, trying to get it from the request header...")
			authToken := c.Get("Authorization")

			if authToken == "" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "unauthorized.",
				})
			}

			authParts := strings.Split(authToken, " ")

			if len(authParts) != 2 || authParts[0] != "Bearer" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "unauthorized.",
				})
			}
			token = authParts[1]
		}

		secret := []byte("super-secret-key")

		jwttoken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
				return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !jwttoken.Valid {
			c.ClearCookie()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized.",
			})
		}

		userId := jwttoken.Claims.(jwt.MapClaims)["userId"]

		if err := db.Model(&User{}).Where("id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			c.ClearCookie()
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized.",
			})
		}
		c.Locals("userid", userId)

		return c.Next()
	}
}
