package main
import(
  
  "context"
  "io/ioutil"
  "net/http"
  "noah/clypeus-dashboard/structures"
  "github.com/joho/godotenv"
  "go/oauth2"
  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/utils"
)
var (
	state     = utils.UUID()
	OauthConf *oauth2.Config
)

var (
	CLIENT_ID     string
	CLIENT_SECRET string
)

	godotenv.Load()
	CLIENT_ID = os.Getenv("CLIENT_ID")
	CLIENT_SECRET = os.Getenv("CLIENT_SECRET")

func declareAuth(authGroup fiber.Router) { 
 authGroup.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/auth/login")
	})
	authGroup.Get("/login", func(c *fiber.Ctx) error {
		OauthConf = &oauth2.Config{
			RedirectURL:  "http://localhost:3000/auth/callback",
			ClientID:     CLIENT_ID,
			ClientSecret: CLIENT_SECRET,
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
}
