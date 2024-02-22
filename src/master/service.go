package master

import (
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type MasterService struct {
	CS              *cluster.ClusterState
	NodesStatus     map[string]cluster.NodeStatus
	masterLog       *cluster.Log
	NodesStatusLogs map[string]*cluster.Log
	LogsPath        string
}

func NewMasterService() *MasterService {
	m := &MasterService{
		NodesStatusLogs: make(map[string]*cluster.Log),
		NodesStatus:     make(map[string]cluster.NodeStatus),
		LogsPath:        "./logs/masterLogs/",
	}
	m.masterLog = cluster.NewLog(m.LogsPath + "masterLog")
	return m
}
