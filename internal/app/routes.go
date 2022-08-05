package app

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type AccountModel struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
	RegionCode  string `json:"regionCode"`
}

type ChangeHistoryEventsRequest struct {
	EarliestChangeTime string   `json:"earliestChangeTime"`
	PageSize           int      `json:"pageSize"`
	ResourceType       []string `json:"resourceType"`
}

type ChangeHistoryEventModel struct {
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

	ch := make(chan []ChangeHistoryEventModel)

	for _, acc := range accounts {
		go app.GetChangeHistory(acc, client, ch)
	}

	app.convertJSONToCSV(<-ch, []string{"ChangeTime", "UserActorEmail", "ActorType"}, "changeHistoryEvents.csv")

}

func (app *Application) ListAccounts(c *http.Client) []AccountModel {

	res, err := c.Get("https://analyticsadmin.googleapis.com/v1beta/accounts/?pageSize=200")

	if err != nil {
		log.Fatalln(err)
	}

	var response struct {
		Accounts []AccountModel `json:"accounts"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)

	if err != nil {
		panic(err)
	}

	var accountsArray []AccountModel

	accountsList := append(accountsArray, response.Accounts...)

	return accountsList

}

func (app *Application) GetChangeHistory(acc AccountModel, c *http.Client, ch chan []ChangeHistoryEventModel) {

	url := fmt.Sprintf("https://analyticsadmin.googleapis.com/v1beta/%s:searchChangeHistoryEvents", acc.Name)

	postBody := &ChangeHistoryEventsRequest{
		EarliestChangeTime: "2022-07-01T00:00:00.000Z",
		PageSize:           1000,
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

	b, err := json.Marshal(postBody)

	if err != nil {
		app.logger.Fatalf("Failed to marshal body: %s", err)
	}

	res, err := c.Post(url, "application/json", bytes.NewBuffer(b))

	if err != nil {
		app.logger.Fatalf("Post request failed: %s", err)
	}

	var result struct {
		ChangeHistoryEvents []ChangeHistoryEventModel `json:"changeHistoryEvents"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		app.logger.Fatalf("Failed to decode response: %s", err)
	}

	var changeHistoryEventModelArray []ChangeHistoryEventModel

	if len(result.ChangeHistoryEvents) > 0 {
		changeHistoryEvents := append(changeHistoryEventModelArray, result.ChangeHistoryEvents...)
		ch <- changeHistoryEvents
	}

}

func (app Application) convertJSONToCSV(c []ChangeHistoryEventModel, header []string, destination string) error {

	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 4. Write the header of the CSV file and the successive rows by iterating through the JSON struct array
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range c {
		var csvRow []string
		csvRow = append(csvRow, r.ChangeTime, r.UserActorEmail, r.ActorType)
		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}

// func (app *Application) GetChangeHistory(acc []AccountModel, c *http.Client) {
// 	f, err := os.OpenFile("change_history", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// 	if err != nil {
// 		app.logger.Fatalf("Failed creating a file: %s", err)
// 	}

// 	w := bufio.NewWriter(f)

// 	var wg sync.WaitGroup
// 	wg.Add(len(acc))

// 	for _, account := range acc {
// 		go func(account AccountModel) {
// 			url := fmt.Sprintf("https://analyticsadmin.googleapis.com/v1beta/%s:searchChangeHistoryEvents", account.Name)

// 			postBody := &ChangeHistoryEventsRequest{
// 				EarliestChangeTime: "2022-07-01T00:00:00.000Z",
// 				PageSize:           1000,
// 				ResourceType: []string{
// 					"ACCOUNT",
// 					"PROPERTY",
// 					"GOOGLE_ADS_LINK",
// 					"GOOGLE_SIGNALS_SETTINGS",
// 					"CONVERSION_EVENT",
// 					"MEASUREMENT_PROTOCOL_SECRET",
// 					"DATA_RETENTION_SETTINGS",
// 					"DISPLAY_VIDEO_360_ADVERTISER_LINK",
// 					"DISPLAY_VIDEO_360_ADVERTISER_LINK_PROPOSAL",
// 					"DATA_STREAM",
// 					"ATTRIBUTION_SETTINGS",
// 				},
// 			}

// 			b, err := json.Marshal(postBody)

// 			if err != nil {
// 				app.logger.Fatalf("Failed to marshal body: %s", err)
// 			}

// 			res, err := c.Post(url, "application/json", bytes.NewBuffer(b))

// 			if err != nil {
// 				app.logger.Fatalf("Post request failed: %s", err)
// 			}

// 			body, err := ioutil.ReadAll(res.Body)

// 			fmt.Println(string(body))

// 			if err != nil {
// 				app.logger.Fatalf("Failed to read Post body: %s", err)
// 			}

// 			var result map[string]interface{}

// 			json.Unmarshal([]byte(body), &result)

// 			if result["error"] != nil {
// 				wg.Done()
// 				return
// 			}

// 			_, err = w.WriteString(string(body))

// 			if err != nil {
// 				app.logger.Fatalln("Error writing to a file:", err)
// 			}
// 			wg.Done()
// 		}(account)

// 	}
// 	wg.Wait()
// 	w.Flush()
// 	f.Close()
// 	os.Exit(0)
// }
