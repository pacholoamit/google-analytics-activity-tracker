package models

type Account struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
	RegionCode  string `json:"regionCode"`
}

type ChangeHistoryEvent struct {
	ChangeTime      string `json:"changeTime"`
	ActorType       string `json:"actorType"`
	UserActorEmail  string `json:"userActorEmail"`
	ChangesFiltered bool   `json:"changesFiltered"`
	Changes         []struct {
		Resource string `json:"resource"`
		Action   string `json:"action"`
	}
}
