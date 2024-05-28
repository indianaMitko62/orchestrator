package master

import (
	"log/slog"
	"os"

	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type MasterSettings struct {
	Name              string `yaml:"name"`
	HTTPServerPort    string `yaml:"httpserver_port"`
	LogsPath          string `yaml:"logs_path"`
	DefaultNetworking bool   `yaml:"use_default_networking"`
}

type MasterService struct {
	MasterSettings
	CS              *cluster.ClusterState
	NodesStatus     map[string]cluster.NodeStatus
	masterLog       *cluster.Log
	NodesStatusLogs map[string]*cluster.Log
	nodeNameToIP    map[string]string
}

func NewMasterService(ms MasterSettings) *MasterService {
	m := &MasterService{
		MasterSettings:  ms,
		NodesStatusLogs: make(map[string]*cluster.Log),
		NodesStatus:     make(map[string]cluster.NodeStatus),
	}
	if err := os.Mkdir(m.LogsPath, 0755); os.IsExist(err) {
		slog.Info("Directory exists", "name", m.LogsPath)
	} else {
		slog.Info("Directory created", "name", m.LogsPath)
	}
	m.CS = cluster.NewClusterState()
	m.masterLog = cluster.NewLog(m.LogsPath + "masterLog")
	m.nodeNameToIP = make(map[string]string)
	return m
}
