package app

import (
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Application struct {
	Config Config
	Oauth  *oauth2.Config
	Logger *log.Logger
}

type Config struct {
	ClientId     string
	ClientSecret string
	RedirectURL  string
}

func New(cfg Config, logger *log.Logger) *Application {
	oauth := &oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/business.manage",
			"https://www.googleapis.com/auth/analytics.readonly",
			"https://www.googleapis.com/auth/analytics.edit",
			"https://www.googleapis.com/auth/adwords",
		},
		Endpoint: google.Endpoint,
	}
	return &Application{
		Config: cfg,
		Oauth:  oauth,
		Logger: logger,
	}
}
