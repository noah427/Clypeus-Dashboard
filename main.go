package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/template/pug"
	"golang.org/x/oauth2"
)

var (
	state     = csrf.New()
	OauthConf *oauth2.Config
)

func main() {
	engine := pug.New("./public/views", ".pug")

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
	authGroup := app.Group("/auth")
	authGroup.Get("/login", func(c *fiber.Ctx) error {
		OauthConf = &oauth2.Config{
			RedirectURL:  "http://localhost:3000/auth/callback",
			ClientID:     "",
			ClientSecret: "",
			Scopes:       []string{"identify", "guilds"},
			Endpoint: oauth2.Endpoint{
				TokenURL: "https://discordapp.com/api/oauth2/token",
				AuthURL:  "https://discordapp.com/api/oauth2/authorize",
			},
		}
		url := OauthConf.AuthCodeURL(state, oauth2.AccessTypeOnline)
		url += "&prompt=none"
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
		Reqbody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			c.SendStatus(http.StatusInternalServerError)
			return nil
		}
		marshalled, err := json.Marshal(Reqbody)
		if err != nil {
			log.Printf("Could not marshal data")
		}
		return c.JSONP(marshalled, "discord-data")
	})
	app.Static("/js", "./public/js")
	app.Static("/fonts", "./public/fonts")
	app.Static("/images", "./public/images")
	app.Static("/css", "./public/css")
	app.Listen(":3000")
}
