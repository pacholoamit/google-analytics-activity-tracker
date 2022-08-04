package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pacholoamit/google-analytics-activity-monitor/internal/app"
	"github.com/pacholoamit/google-analytics-activity-monitor/internal/config"
	"github.com/spf13/viper"
)

func main() {
	conf := setupConfig()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := app.New(conf, logger)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 3000),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Printf("starting server on %s", srv.Addr)
		err := srv.ListenAndServe()
		logger.Fatal(err)
	}()
	app.GoogleAuthenticate()
	app.ListAccounts()

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
