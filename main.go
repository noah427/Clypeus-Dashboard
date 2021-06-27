package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/pug"
)


func main() {
	engine := pug.New("./public/views", ".pug")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	api := app.Group("/api")

	declareApi(api)

	//load env
	godotenv.Load()
	CLIENT_ID = os.Getenv("CLIENT_ID")
	CLIENT_SECRET = os.Getenv("CLIENT_SECRET")

	loadDatabase()

	authGroup := app.Group("/auth")
	declareAuth(authGroup)
	dashboard := app.Group("/dashboard")
	dashboard.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/dashboard/selector")
	})
	dashboard.Get("/selector", func(c *fiber.Ctx) error {
		return c.Render("selector", fiber.Map{})
	})
	app.Static("/js", "./public/js")
	app.Static("/fonts", "./public/fonts")
	app.Static("/images", "./public/images")
	app.Static("/css", "./public/css")
	app.Listen(":3000")
}
