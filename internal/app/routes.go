package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type accountModel struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
	RegionCode  string `json:"regionCode"`
}

type changeHistoryEventModel struct {
	ChangeTime      string `json:"changeTime"`
	ActorType       string `json:"actorType"`
	UserActorEmail  string `json:"userActorEmail"`
	ChangesFiltered bool   `json:"changesFiltered"`
	Changes         []struct {
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}
}

func (app *Application) Routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/success", app.successHandler)
	return router
}

func (app *Application) successHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	app.logger.Println("Code successfully requested:", code)

	client := app.newGoogleClient(code)

	accounts := app.ListAccounts(client)

	app.GetChangeHistory(accounts, client)
}

func (app *Application) ListAccounts(c *http.Client) []accountModel {

	res, err := c.Get("https://analyticsadmin.googleapis.com/v1beta/accounts/?pageSize=200")

	if err != nil {
		log.Fatalln(err)
	}

	var response struct {
		Accounts []accountModel `json:"accounts"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		panic(err)
	}

	var accountsArray []accountModel

	accountsList := append(accountsArray, response.Accounts...)

	return accountsList

}

func (app *Application) GetChangeHistory(acc []accountModel, c *http.Client) {
	f, err := os.OpenFile("change_history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		app.logger.Fatalf("Failed creating a file: %s", err)
	}

	w := bufio.NewWriter(f)

	var wg sync.WaitGroup

	wg.Add(len(acc))

	for _, account := range acc {
		go func(account accountModel) {
			url := fmt.Sprintf("https://analyticsadmin.googleapis.com/v1beta/%s:searchChangeHistoryEvents", account.Name)

			type ChangeHistoryEventsRequest struct {
				EarliestChangeTime string   `json:"earliestChangeTime"`
				PageSize           int      `json:"pageSize"`
				ResourceType       []string `json:"resourceType"`
			}

			postBody := &ChangeHistoryEventsRequest{
				EarliestChangeTime: "2022-07-01T00:00:00.000Z",
				PageSize:           100,
				ResourceType: []string{
					"ACCOUNT",
					"PROPERTY",
					"GOOGLE_ADS_LINK",
					"GOOGLE_SIGNALS_SETTINGS",
					"CONVERSION_EVENT",
					"MEASUREMENT_PROTOCOL_SECRET",
					"DATA_RETENTION_SETTINGS",
					"DISPLAY_VIDEO_360_ADVERTISER_LINK",
					"DISPLAY_VIDEO_360_ADVERTISER_LINK_PROPOSAL",
					"DATA_STREAM",
					"ATTRIBUTION_SETTINGS",
				},
			}

			b, _ := json.Marshal(postBody)

			res, err := c.Post(url, "application/json", bytes.NewBuffer(b))

			if err != nil {
				app.logger.Fatalln(err)
			}

			body, err := ioutil.ReadAll(res.Body)

			fmt.Println(string(body))

			if err != nil {
				app.logger.Fatalln(err)
			}

			var result map[string]interface{}

			json.Unmarshal([]byte(body), &result)

			if result["error"] != nil {
				wg.Done()
				return
			}

			_, err = w.WriteString(string(body))

			if err != nil {
				app.logger.Fatalln("Error writing to a file:", err)
			}
			wg.Done()
		}(account)

	}
	wg.Wait()
	w.Flush()
	f.Close()
	os.Exit(0)
}
