package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"noah/clypeus-dashboard/structures"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/template/pug"
	"golang.org/x/oauth2"
)

var (
	state     = utils.UUID()
	OauthConf *oauth2.Config
)

func main() {
	engine := pug.New("./public/views", ".pug")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	authGroup := app.Group("/auth")
	authGroup.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/auth/login")
	})
	authGroup.Get("/login", func(c *fiber.Ctx) error {
		OauthConf = &oauth2.Config{
			RedirectURL:  "http://localhost:3000/auth/callback",
			ClientID:     "781817418535403531",
			ClientSecret: "MmY_uTCQ1iQJcNwMn8osdpu_WfSzQn8f",
			Scopes:       []string{"identify", "guilds"},
			Endpoint: oauth2.Endpoint{
				TokenURL: "https://discordapp.com/api/oauth2/token",
				AuthURL:  "https://discordapp.com/api/oauth2/authorize",
			},
		}
		url := OauthConf.AuthCodeURL(state, oauth2.AccessTypeOnline)
		return c.Redirect(url, http.StatusTemporaryRedirect)
	})
	authGroup.Get("/callback", func(c *fiber.Ctx) error {
		formValue := c.FormValue("state")
		if formValue != state {
			c.SendStatus(http.StatusBadRequest)
			return nil
		}
		token, err := OauthConf.Exchange(context.Background(), c.FormValue("code"))
		if err != nil {
			c.SendStatus(http.StatusBadRequest)

			return nil
		}
		res, err := OauthConf.Client(context.Background(), token).Get("https://discordapp.com/api/users/@me")
		if err != nil {
			c.SendStatus(http.StatusInternalServerError)
			return nil
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		err = json.Unmarshal(body, &structures.Data)
		if err != nil {
			c.SendStatus(http.StatusBadRequest)
		}
		return c.Redirect("/dashboard")
	})
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
