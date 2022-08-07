package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pacholoamit/google-analytics-activity-monitor/internal/app"
	"github.com/pacholoamit/google-analytics-activity-monitor/internal/client"
)

func main() {
	var cfg client.Config
	var appCfg app.Config

	flag.StringVar(&cfg.ClientId, "clientId", "", "Google Client ID")
	flag.StringVar(&cfg.ClientSecret, "clientSecret", "", "Google Client Secret")
	flag.StringVar(&appCfg.File, "file", "", "File output path")
	flag.Parse()

	if err := cfg.ValidateFlags(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := appCfg.ValidateFlags(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	c := client.New(cfg, l)

	app := app.New(c, l, &appCfg)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 3000),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Client.Authenticate()

	app.Logger.Printf("starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	app.Logger.Fatal(err)

}
