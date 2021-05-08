package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/pug"
)

func main() {

	engine := pug.New("./public/views", ".pug")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
		})
	})

	app.Static("/js", "./public/js")
	app.Static("/fonts", "./public/fonts")
	app.Static("/images", "./public/images")
	app.Static("/css", "./public/css")
	app.Listen(":3000")
}
