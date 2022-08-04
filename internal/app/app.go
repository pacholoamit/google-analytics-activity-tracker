package app

import (
	"io/ioutil"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Application struct {
	config config
	oauth  *oauth2.Config
	logger *log.Logger
}

type config interface {
	GetClientId() string
	GetClientSecret() string
	GetRedirectURL() string
	GetCode() string
	SetCode(code string)
}

func New(cfg config, logger *log.Logger) *Application {
	oauth := &oauth2.Config{
		ClientID:     cfg.GetClientId(),
		ClientSecret: cfg.GetClientSecret(),
		RedirectURL:  cfg.GetRedirectURL(),
		Scopes: []string{
			"https://www.googleapis.com/auth/business.manage",
			"https://www.googleapis.com/auth/analytics.readonly",
			"https://www.googleapis.com/auth/adwords",
		},
		Endpoint: google.Endpoint,
	}
	return &Application{
		config: cfg,
		oauth:  oauth,
		logger: logger,
	}
}

func (app *Application) ListAccounts() {
	client := app.newGoogleClient()

	resp, err := client.Get("https://analyticsadmin.googleapis.com/v1alpha/accounts")

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)

	log.Print(sb)
}
