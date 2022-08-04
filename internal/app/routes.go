package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/success", app.successHandler)
	return router
}

func (app *Application) successHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	app.config.SetCode(code)
	app.logger.Println("Code successfully requested:", code)
	app.ListAccounts()
}
