package app

import (
	"context"
	"fmt"
	"net/http"
)

func (app *Application) GoogleAuthenticate() {
	url := app.Oauth.AuthCodeURL("state") // For inclusing of refresh token
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
}

func (app *Application) newGoogleClient(code string) *http.Client {
	token, err := app.Oauth.Exchange(context.Background(), code)
	if err != nil {
		app.Logger.Fatal(err)
	}
	return app.Oauth.Client(context.Background(), token)
}

func (c *Config) ValidateFlags() error {
	if c.ClientId == "" || c.ClientSecret == "" || c.RedirectURL == "" {
		return fmt.Errorf(`please provide all required flags:
			-clientId
				Google Client ID
			-clientSecret
				Google Client Secret
			-redirectUrl
				Google authorized redirect URL
			`)
	}
	return nil
}
