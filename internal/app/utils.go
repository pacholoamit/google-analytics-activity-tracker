package app

import (
	"context"
	"net/http"
)

func (app *Application) newGoogleClient(code string) *http.Client {
	token, err := app.oauth.Exchange(context.Background(), code)
	if err != nil {
		app.logger.Fatal(err)
	}
	return app.oauth.Client(context.Background(), token)
}
