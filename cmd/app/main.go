package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Enter you Client Id: ")
		var clientId string
		fmt.Scanln(&clientId)
		viper.Set("client_id", clientId)

		fmt.Println("Enter you Client Secret: ")
		var clientSecret string
		fmt.Scanln(&clientSecret)
		viper.Set("client_secret", clientSecret)

		fmt.Println("Enter you Redirect URL: ")
		var redirectURL string
		fmt.Scanln(&redirectURL)
		viper.Set("redirect_url", redirectURL)

		viper.SafeWriteConfig()
	}

	clientId := viper.GetString("client_id")
	clientSecret := viper.GetString("client_secret")
	redirectURL := viper.GetString("redirect_url")

	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/business.manage",
			"https://www.googleapis.com/auth/analytics.readonly",
			"https://www.googleapis.com/auth/adwords",
		},
		Endpoint: google.Endpoint,
	}

	url := conf.AuthCodeURL("state") // For inclusing of refresh token

	browser.OpenURL(url)

}

type config struct {
	clientId     string
	clientSecret string
	redirectURL  string
}

type application struct {
	config config
	oauth  *oauth2.Config
}

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/success", app.successHandler)
	return router
}

func (app *application) successHandler(w http.ResponseWriter, r *http.Request) {

	token, err := app.oauth.Exchange(context.Background(), "authorization-code")

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
