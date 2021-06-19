package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func declareApi(api fiber.Router) {
	api.Post("/guild/:id/settings", func(c *fiber.Ctx) error {
		id := c.Params("id")
		settings := &GuildSettings{}
		c.BodyParser(settings)
		fmt.Println(string(c.Body()))
		settings.ID = id
		fmt.Println(settings.AntinukeON)
		return nil
	})
}
