package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	// Initialize the DB
	db := InitializeDB()

	// Create the Fiber App Instance
	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})

	// Auth routes
	AuthHandlers(app.Group("/auth"), db)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
