package model

type ServiceInfo struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	Clustername string `json:"clustername"`
	Status string `json:"status"`
}

type ServiceInfos []ServiceInfo
