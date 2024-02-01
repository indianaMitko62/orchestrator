package master

import (
	"github.com/docker/docker/client"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type MasterService struct {
	cli *client.Client

	CS *cluster.ClusterState
}

func NewMasterService(cs *cluster.ClusterState) *MasterService {
	m := &MasterService{}
	m.CS = cs
	return m
}
