package master

import (
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type MasterService struct {
	CS          *cluster.ClusterState
	NodesStatus map[string]cluster.NodeStatus
	MasterLog   *cluster.Log
}

func NewMasterService() *MasterService {
	m := &MasterService{
		MasterLog:   cluster.NewLog("./logs/masterLog"),
		NodesStatus: make(map[string]cluster.NodeStatus),
	}
	return m
}
