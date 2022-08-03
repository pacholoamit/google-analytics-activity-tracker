package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pkg/browser"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

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

	token, err := conf.Exchange(context.Background(), "authorization-code")

	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(context.Background(), token)

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
