package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pacholoamit/google-analytics-activity-monitor/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Google struct {
	config Config
	oauth  *oauth2.Config
	logger *log.Logger
}

type Config struct {
	ClientId     string
	ClientSecret string
}

type ChangeHistoryEventsRequest struct {
	EarliestChangeTime string   `json:"earliestChangeTime"`
	PageSize           int      `json:"pageSize"`
	ResourceType       []string `json:"resourceType"`
}

func New(cfg Config, l *log.Logger) *Google {
	oauth := &oauth2.Config{
		ClientID:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  "http://localhost:3000/success",
		Scopes: []string{
			"https://www.googleapis.com/auth/business.manage",
			"https://www.googleapis.com/auth/analytics.readonly",
			"https://www.googleapis.com/auth/analytics.edit",
			"https://www.googleapis.com/auth/adwords",
		},
		Endpoint: google.Endpoint,
	}
	return &Google{
		oauth:  oauth,
		config: cfg,
		logger: l,
	}
}

func (c *Google) Authenticate() {
	url := c.oauth.AuthCodeURL("state") // For inclusing of refresh token
	c.logger.Printf("Visit the URL for the auth dialog: %v", url)
}

func (c *Google) Exchange(code string) *http.Client {
	token, err := c.oauth.Exchange(context.Background(), code)
	if err != nil {
		c.logger.Fatal(err)
	}
	return c.oauth.Client(context.Background(), token)
}

func (c *Google) ListAccounts(http *http.Client) []models.Account {
	res, err := http.Get("https://analyticsadmin.googleapis.com/v1beta/accounts/?pageSize=200")

	if err != nil {
		log.Fatalln(err)
	}

	var response struct {
		Accounts []models.Account `json:"accounts"`
	}

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		log.Fatalln(err)
	}

	var acc []models.Account

	accounts := append(acc, response.Accounts...)

	return accounts
}

func (c *Google) GetChangeHistory(http *http.Client, accountName string, ch chan []models.ChangeHistoryEvent) {
	url := fmt.Sprintf("https://analyticsadmin.googleapis.com/v1beta/%s:searchChangeHistoryEvents", accountName)

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
		c.logger.Fatalf("Failed to marshal body: %s", err)
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(b))

	if err != nil {
		c.logger.Fatalf("Post request failed: %s", err)
	}

	var response struct {
		ChangeHistoryEvents []models.ChangeHistoryEvent `json:"changeHistoryEvents"`
	}

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		c.logger.Fatalf("Failed to decode response: %s", err)
	}

	var events []models.ChangeHistoryEvent

	if len(response.ChangeHistoryEvents) > 0 {
		ch <- append(events, response.ChangeHistoryEvents...)
	}

}
