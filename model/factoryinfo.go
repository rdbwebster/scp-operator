package model

import (
	api "github.com/rdbwebster/scp-operator/api/v1"
)

// No longer used TODO remove
type FactoryInfo struct {
	Spec        api.ManagedOperatorSpec `json:"spec"`
	Clustername string                  `json:"clustername"`
}

type FactoryInfos []FactoryInfo
