package models

type Account struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
	RegionCode  string `json:"regionCode"`
}

type ChangeHistoryEvent struct {
	ChangeTime      string    `json:"changeTime"`
	ActorType       string    `json:"actorType"`
	UserActorEmail  string    `json:"userActorEmail"`
	ChangesFiltered bool      `json:"changesFiltered"`
	Changes         []Changes `json:"changes"`
}

type Changes struct {
	Resource             string                `json:"resource"`
	Action               string                `json:"action"`
	ResourceBeforeChange ChangeHistoryResource `json:"resourceBeforeChange"`
	ResourceAfterChange  ChangeHistoryResource `json:"resourceAfterChange"`
}

type ChangeHistoryResource struct {
	Account ChangeHistoryAccount `json:"account"`
}

type ChangeHistoryAccount struct {
	Name        string `json:"name"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
	DisplayName string `json:"displayName"`
	RegionCode  string `json:"regionCode"`
	Deleted     bool   `json:"deleted"`
}
