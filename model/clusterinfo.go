package model

import (
	api "github.com/rdbwebster/scp-operator/api/v1"
)

type ClusterInfo struct {
	Spec api.SCPclusterSpec `json:"spec"`
}

type ClusterInfos []ClusterInfo
