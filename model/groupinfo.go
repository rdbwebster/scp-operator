package model

type GroupInfo struct {
	Name   string `json:"name"`
	members    []string `json:"member"`
}

type GroupInfos []GroupInfo
