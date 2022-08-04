package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
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

	fmt.Print("Code: ", code)
	token, err := app.oauth.Exchange(context.Background(), code)

	if err != nil {
		log.Fatal(err)
	}

	client := app.oauth.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/analytics/v3/management/accounts")

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
