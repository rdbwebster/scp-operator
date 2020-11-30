package model

type FactoryInfo struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	Status string `json:"status"`
}

type FactoryInfos []FactoryInfo
