package app

import (
	"log"
	"net/http"
	"os"
	"sync"

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
	GetChangeHistory(h *http.Client, accountName string) ([]models.ChangeHistoryEvent, error)
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

	ch := []models.ChangeHistoryEvent{}
	wg := sync.WaitGroup{}
	// mu := &sync.Mutex{}

	wg.Add(len(accounts))
	for i, account := range accounts {
		go func(i int, account models.Account) {
			hist, err := app.Client.GetChangeHistory(client, account.Name)
			if err != nil {
				app.Logger.Fatalln(err)
			}
			ch = append(ch, hist...)

			wg.Done()
		}(i, account)
	}

	wg.Wait()
	// b, err := json.Marshal(ch)

	// if err != nil {
	// 	app.Logger.Fatalln(err)
	// }
	// var jsonData bytes.Buffer

	// json.Indent(&jsonData, b, "", "  ")

	// app.Logger.Println(jsonData.String())
	app.Logger.Print("Channel closed :", len(ch))

	headers := []string{"UserActorEmail", "ChangeTime", "ActorType", "Changes"}

	if err := app.writeJSONToCSV(ch, headers, "changeHistoryEvents.csv"); err != nil {
		app.Logger.Fatalln("Failed to write JSON to csv: ", err)
	}

	os.Exit(0)
}
