package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pacholoamit/google-analytics-activity-monitor/internal/app"
)

func main() {
	var c app.Config

	flag.StringVar(&c.ClientId, "clientId", "", "Google Client ID")
	flag.StringVar(&c.ClientSecret, "clientSecret", "", "Google Client Secret")
	flag.Parse()

	if err := c.ValidateFlags(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := app.New(c, l)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 3000),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.GoogleAuthenticate()

	app.Logger.Printf("starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	app.Logger.Fatal(err)

}
