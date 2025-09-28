package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DownloadHandler(route fiber.Router, db *gorm.DB) {
	route.Get("/", func(c *fiber.Ctx) error {
		userId := c.Locals("userId").(int)
		books := new([]Book)

		format := c.Query("format", "json")

		if err := db.Where("userId = ?", userId).Find(&books).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var fileName string

		switch format {
		case "json":
			fileName = "books.json"
			file, err := os.Create(fileName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create a JSON file",
				})
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", " ")

			if err := encoder.Encode(books); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to write into JSON file",
				})
			}

		case "csv":
			fileName = "books.csv"
			file, err := os.Create(fileName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to create a CSV file",
				})
			}
			defer file.Close()

			writer := csv.NewWriter(file)

			writer.Write([]string{"ID", "Title", "Status", "Author", "Year", "userId"})

			for _, book := range *books {
				writer.Write([]string{
					fmt.Sprintf("%v", book.ID),
					book.Title,
					string(book.Status),
					book.Author,
					fmt.Sprintf("%v", book.UserID),
				})
			}

			writer.Flush()

			if err := writer.Error(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to write into CSV file",
				})
			}

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Only supports json and csv",
			})
		}

		defer os.Remove(fileName)
		return c.Download("./" + fileName)
	})
}
