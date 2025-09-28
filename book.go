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

	route.Get("/", func(c *fiber.Ctx) error {
		userId := c.Locals("userId").(int)
		books := new([]Book)

		title := c.Query("title")
		status := c.Query("status")
		author := c.Query("author")
		year := c.QueryInt("year")

		query := db.Where("userId = ?", userId)

		if title != "" {
			query.Where("title LIKE ?", "%"+title+"%")
		}

		if status != "" {
			query.Where("status = ?", status)
		}

		if author != "" {
			query.Where("author = ?", author)
		}

		if year != 0 {
			query.Where("year = ?", year)
		}

		if err := query.Find(&books).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "book not found.",
			})
		}
		return c.Status(fiber.StatusOK).JSON(books)
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
				"error": "book not found.",
			})
		}
		return c.Status(fiber.StatusOK).JSON(book)
	})

	route.Put("/:id", func(c *fiber.Ctx) error {
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
				"error": "book not found.",
			})
		}

		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := db.Save(&book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(book)
	})

	route.Delete("/:id", func(c *fiber.Ctx) error {
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
				"error": "book not found.",
			})
		}

		if err := db.Delete(&book).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
}
