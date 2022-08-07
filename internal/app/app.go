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
	Config *Config
	Logger *log.Logger
}

type GoogleClient interface {
	Authenticate()
	Exchange(code string) *http.Client
	ListAccounts(h *http.Client) []models.Account
	GetChangeHistory(h *http.Client, accountName string) ([]models.ChangeHistoryEvent, error)
}

type envelope map[string]interface{}

func New(c GoogleClient, l *log.Logger, cfg *Config) *Application {
	return &Application{
		Client: c,
		Config: cfg,
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

	app.Logger.Print("Channel closed :", len(ch))

	a := envelope{"activity": ch}

	if err := app.writeJSONToFile(app.Config.File, a); err != nil {
		app.Logger.Fatalln("Failed to marshal json file: ", err)
	}

	app.Logger.Print("File written to: ", app.Config.File)

	os.Exit(0)
}
