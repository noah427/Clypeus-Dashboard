package main

import (
	"github.com/gofiber/fiber/v2"
)

func declareApi(api fiber.Router) {
	api.Post("/guild/:id/settings", func(c *fiber.Ctx) error {
		var err error
		id := c.Params("id")
		settings := &GuildSettings{}
		err = c.BodyParser(settings)
		if err != nil {
			return err
		}
		settings.ID = id
		db.Create(&settings)
		return nil
	})

	api.Get("/guild/:id/settings", func(c *fiber.Ctx) error {
		var err error
		id := c.Params("id")
		settings := &GuildSettings{ID: id}
		db.Take(&settings)
		err = c.JSON(settings)
		if err != nil {
			return err
		}
		return nil
	})
}
