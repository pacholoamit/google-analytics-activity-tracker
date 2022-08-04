package app

import (
	"log"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
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
}

func New(cfg config, oauth *oauth2.Config, logger *log.Logger) *Application {
	return &Application{
		config: cfg,
		oauth:  oauth,
		logger: logger,
	}
}

func (app *Application) GoogleAuthenticate() {
	url := app.oauth.AuthCodeURL("state") // For inclusing of refresh token
	browser.OpenURL(url)
}
