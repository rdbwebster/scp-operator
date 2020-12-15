package model

type ServiceInfo struct {
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	Clustername string    `json:"clustername"`
	Status      string    `json:"status"`
	CRoutputs   []CRentry `json:"crvalues,omitempty"`
}

type CRentry struct {
	ControlName  string `json:"controlname"`
	ValueType    string `json:"valuetype"`
	ControlValue string `json:"controlvalue"`
	CRpath       string `json:"crpath"`
}

type ServiceInfos []ServiceInfo
