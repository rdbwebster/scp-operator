package model

type GroupInfo struct {
	Name    string   `json:"name"`
	Members []string `json:"member"`
}

type GroupInfos []GroupInfo
