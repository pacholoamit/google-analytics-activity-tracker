package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
)

func (app *Application) Routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/success", app.successHandler)
	return router
}

func (app *Application) successHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	prevConfig := viper.AllSettings()

	for k, v := range prevConfig {
		viper.Set(k, v)
	}
	viper.Set("code", code)
	err := viper.WriteConfig()

	if err != nil {
		app.logger.Printf("Error writing config: %s", err)
	}

	// client := app.newGoogleClient(code)

	// resp, err := client.Get("https://analyticsadmin.googleapis.com/v1alpha/accounts")

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// body, err := ioutil.ReadAll(resp.Body)

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// sb := string(body)

	// log.Print(sb)
}
