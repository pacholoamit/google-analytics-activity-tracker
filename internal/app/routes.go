package app

import (
	"encoding/json"
	"log"
	"net/http"

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
	app.ListAccounts()
}

type account struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
	RegionCode  string `json:"regionCode"`
}

func (app *Application) ListAccounts() {
	client := app.newGoogleClient()

	resp, err := client.Get("https://analyticsadmin.googleapis.com/v1alpha/accounts/?pageSize=200")

	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	var responseJSON struct {
		Accounts []account `json:"accounts"`
	}

	err = json.NewDecoder(resp.Body).Decode(&responseJSON)

	if err != nil {
		panic(err)
	}

	for _, account := range responseJSON.Accounts {

		log.Println(account.Name)
	}
	// body, err := ioutil.ReadAll(resp.Body)

	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// sb := string(body)

	// log.Print(sb)
}
