package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/success", app.successHandler)
	return router
}

func (app *Application) successHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	app.config.SetCode(code)

	app.logger.Println("Code successfully requested:", code)
	client := app.newGoogleClient()

	accounts := app.ListAccounts(client)

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
	f, err := os.Create("changeHistory")
	if err != nil {
		app.logger.Fatal(err)
	}
	for _, account := range acc {
		go func(account accountModel) {
			url := fmt.Sprintf("https://analyticsadmin.googleapis.com/v1beta/%s:searchChangeHistoryEvents", account.Name)

			postBody := []byte(`{
  "earliestChangeTime": "2014-10-02T15:01:23Z",
  "property": "properties/250119597",
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

			if err != nil {
				app.logger.Fatalln(err)
			}

			var result map[string]interface{}

			json.Unmarshal([]byte(body), &result)

			if result["error"] != nil {
				return
			}
			fmt.Print(account.Name, ": ", len(result["changeHistoryEvents"].([]interface{})))
			f.WriteString(fmt.Sprintf("%s\n", account.Name))

		}(account)

	}

}
