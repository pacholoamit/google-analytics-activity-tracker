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

	for _, account := range accounts {
		fmt.Println(account.Name)
	}

	app.GetChangeHistory(accounts, client)
}

type accountModel struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
	RegionCode  string `json:"regionCode"`
}

func (app *Application) ListAccounts(c *http.Client) []accountModel {

	resp, err := c.Get("https://analyticsadmin.googleapis.com/v1beta/accounts/?pageSize=200")

	if err != nil {
		log.Fatalln(err)
	}

	var responseJSON struct {
		Accounts []accountModel `json:"accounts"`
	}

	err = json.NewDecoder(resp.Body).Decode(&responseJSON)

	if err != nil {
		panic(err)
	}

	var accountsArray []accountModel

	accountsList := append(accountsArray, responseJSON.Accounts...)

	return accountsList

}

func (app *Application) GetChangeHistory(acc []accountModel, c *http.Client) {
	fmt.Println("Getting change history for accounts:")
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

			postBody := []byte(`{
  "earliestChangeTime": "2014-10-02T15:01:23Z",
  "pageSize": 1000,
  "resourceType": [
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
    "ATTRIBUTION_SETTINGS"
  ]
}`)
			resp, err := c.Post(url, "application/json", bytes.NewBuffer(postBody))

			if err != nil {
				app.logger.Fatalln(err)
			}

			body, err := ioutil.ReadAll(resp.Body)

			fmt.Println(string(body))

			if err != nil {
				app.logger.Fatalln(err)
			}

			var result map[string]interface{}

			json.Unmarshal([]byte(body), &result)

			if result["error"] != nil {
				return
			}

			fmt.Println("Account: ", account.Name)

			b, err := w.WriteString(string(body))

			if err != nil {
				app.logger.Fatalln("Error writing to a file:", err)
			}

			fmt.Println("Bytes written: ", b)
			wg.Done()
		}(account)

	}
	// defer w.Flush()
	// defer f.Close()
}
