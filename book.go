package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func BookHandlers(route fiber.Router, db *gorm.DB) {
	route.Post("/", func(c *fiber.Ctx) error {
		book := new(Book)
		book.UserID = c.Locals("userId").(int)

		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := db.Create(&book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(book)
	})

	route.Get("/:id", func(c *fiber.Ctx) error {
		bookId, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		userId := c.Locals("userId").(int)
		book := new(Book)

		if err := db.Where("id = ? AND userId = ?", bookId, userId).First(&book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(book)
	})
}
