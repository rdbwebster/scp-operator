package model

import (
	api "github.com/rdbwebster/scp-operator/api/v1"
)

type ServiceInfo struct {
	Name        string        `json:"name"`
	Crdname     string        `json:"crdname"`
	Clustername string        `json:"clustername"`
	Status      string        `json:"status"`
	CRinputs    []api.CRentry `json:"crinputs"`
	CRoutputs   []api.CRentry `json:"croutputs"`
}

type ServiceInfos []ServiceInfo
