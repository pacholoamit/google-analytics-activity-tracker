package main

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/pacholoamit/google-analytics-activity-monitor/internal/config"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type application struct {
	config *config.Config
	oauth  *oauth2.Config
	logger *log.Logger
}

func main() {

	conf := setupConfig()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	oauth := createOAuthConfig(conf)

	app := &application{
		config: conf,
		oauth:  oauth,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 3000),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("starting server on %s", srv.Addr)
	go func() {
		err := srv.ListenAndServe()
		logger.Fatal(err)
	}()

	url := app.oauth.AuthCodeURL("state") // For inclusing of refresh token
	browser.OpenURL(url)

}

func setupConfig() *config.Config {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter you Client Id: ")
		clientId, _ := reader.ReadString('\n')
		clientId = strings.TrimSpace(clientId)
		viper.Set("client_id", clientId)

		fmt.Println("Enter you Client Secret: ")
		clientSecret, _ := reader.ReadString('\n')
		clientSecret = strings.TrimSpace(clientSecret)
		viper.Set("client_secret", clientSecret)

		fmt.Println("Enter you Redirect URL: ")
		redirectURL, _ := reader.ReadString('\n')
		redirectURL = strings.TrimSpace(redirectURL)
		viper.Set("redirect_url", redirectURL)

		viper.SafeWriteConfig()
	}

	clientId := viper.GetString("client_id")
	clientSecret := viper.GetString("client_secret")
	redirectURL := viper.GetString("redirect_url")

	conf := config.New(clientId, clientSecret, redirectURL)
	return conf
}

func createOAuthConfig(c *config.Config) *oauth2.Config {
	oauth := &oauth2.Config{
		ClientID:     c.GetClientId(),
		ClientSecret: c.GetClientSecret(),
		RedirectURL:  c.GetRedirectURL(),
		Scopes: []string{
			"https://www.googleapis.com/auth/business.manage",
			"https://www.googleapis.com/auth/analytics.readonly",
			"https://www.googleapis.com/auth/adwords",
		},
		Endpoint: google.Endpoint,
	}

	return oauth
}

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/success", app.successHandler)
	return router
}

func (app *application) successHandler(w http.ResponseWriter, r *http.Request) {
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
