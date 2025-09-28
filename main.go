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

	// To protect routes, create a auth middleware
	protected := app.Use(AuthMiddleware(db))

	// Book routes, these routes are protected. Requires a valid JWT
	BookHandlers(protected.Group("/book"), db)

	// Downloads the books as either CSV format or JSON format
	DownloadHandler(protected.Group("/download"), db)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
