package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Library API",
	})

	log.Fatal(app.Listen(":3000"))
}
