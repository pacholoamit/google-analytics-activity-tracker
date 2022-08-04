package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/browser"
)

func (app *Application) GoogleAuthenticate() {
	url := app.oauth.AuthCodeURL("state") // For inclusing of refresh token
	browser.OpenURL(url)
}

func (app *Application) newGoogleClient() *http.Client {
	code := app.config.GetCode()
	fmt.Println(code)
	token, err := app.oauth.Exchange(context.Background(), code)
	if err != nil {
		app.logger.Fatal(err)
	}
	return app.oauth.Client(context.Background(), token)
}
