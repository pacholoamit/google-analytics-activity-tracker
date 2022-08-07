package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/pacholoamit/google-analytics-activity-monitor/internal/models"
)

type Application struct {
	Client GoogleClient
	Logger *log.Logger
}

type GoogleClient interface {
	Authenticate()
	Exchange(code string) *http.Client
	ListAccounts(h *http.Client) []models.Account
	GetChangeHistory(h *http.Client, accountName string, ch chan []models.ChangeHistoryEvent)
}

func New(c GoogleClient, l *log.Logger) *Application {
	return &Application{
		Client: c,
		Logger: l,
	}
}

func (app *Application) Routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/success", app.successHandler)
	return router
}

func (app *Application) successHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	app.Logger.Println("Code successfully requested:", code)

	client := app.Client.Exchange(code)

	accounts := app.Client.ListAccounts(client)

	app.Logger.Println("Accounts retrieved: ", len(accounts))

	ch := make(chan []models.ChangeHistoryEvent)

	for _, account := range accounts {
		go app.Client.GetChangeHistory(client, account.Name, ch)
	}

	for range accounts {
		changes := <-ch
		app.Logger.Println("Changes retrieved: ", len(changes))
	}

	headers := []string{"UserActorEmail", "ChangeTime", "ActorType", "Changes"}

	app.writeJSONToCSV(<-ch, headers, "changeHistoryEvents.csv")

	fmt.Println("Written accounts to csv: ", len(<-ch))

	os.Exit(0)
}
