package master

import (
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type MasterService struct {
	CS              *cluster.ClusterState
	NodesStatus     map[string]cluster.NodeStatus
	MasterLog       *cluster.Log
	NodesStatusLogs map[string]*cluster.Log
}

func NewMasterService() *MasterService {
	m := &MasterService{
		MasterLog:       cluster.NewLog("./logs/masterLogs/masterLog"),
		NodesStatusLogs: make(map[string]*cluster.Log),
		NodesStatus:     make(map[string]cluster.NodeStatus),
	}
	return m
}
