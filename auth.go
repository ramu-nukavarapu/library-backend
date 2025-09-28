package main

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AuthHandlers(route fiber.Router, db *gorm.DB) {
	route.Post("/register", func(c *fiber.Ctx) error {
		user := &User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if user.Username == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "username and password required",
			})
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		user.Password = string(hashed)

		db.Create(user)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User registered successfully.",
		})
	})
	route.Post("/login", func(c *fiber.Ctx) error {
		AuthUser := &User{
			Username: c.FormValue("username"),
			Password: c.FormValue("password"),
		}

		if AuthUser.Username == "" || AuthUser.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "username and password required",
			})
		}

		DbUser := new(User)

		db.Where("username = ?", AuthUser.Username).First(DbUser)

		if DbUser.ID == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "username not found",
			})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(DbUser.Password), []byte(AuthUser.Password)); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid credentials.",
			})
		}

		token, err := GenerateToken(DbUser)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			HTTPOnly: !c.IsFromLocal(),
			Secure:   !c.IsFromLocal(),
			MaxAge:   3600 * 24 * 7,
		})

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	})
}
