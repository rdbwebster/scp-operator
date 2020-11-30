package model

import "time"

type ClusterInfo struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	Token     string    `json:"token"`
	Cert      string    `json:"cert"`
	CertAuth  string    `json:"certauth"`
	Connected time.Time `json:"connected"`
}

type ClusterInfos []ClusterInfo
