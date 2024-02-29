package node

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/docker/docker/client"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type NodeSettings struct {
	Name                 string `yaml:"name"`
	MasterAddress        string `yaml:"master_address"`
	MasterPort           string `yaml:"master_port"`
	LogsEndpoint         string `yaml:"logs_endpoint"`
	ClusterStateEndpoint string `yaml:"clusterstate_endpoint"`
	NodeStatusEndpoint   string `yaml:"nodestatus_endpoint"`
}

type NodeService struct {
	NodeSettings
	cli              *client.Client
	DesiredNodeState *cluster.NodeState // no pointer required
	CurrentNodeState *cluster.NodeState
	clusterChangeLog *cluster.Log
	nodeLog          *cluster.Log
	LogsDir          string
}

func NewNodeService(nodeSetting NodeSettings) (*NodeService, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	ns := &NodeService{
		cli:          cli,
		NodeSettings: nodeSetting,
	}
	ns.LogsDir = "./logs/" + ns.Name + "Logs/"
	if err := os.Mkdir(ns.LogsDir, 0755); os.IsExist(err) {
		slog.Info("Directory exists", "name", ns.LogsDir)
	} else {
		slog.Info("Directory created", "name", ns.LogsDir)
	}
	fmt.Println(ns.NodeSettings.ClusterStateEndpoint)
	ns.clusterChangeLog = cluster.NewLog(ns.LogsDir + "clusterChangeLog")
	ns.nodeLog = cluster.NewLog(ns.LogsDir + "nodeLog")
	return ns, nil
}
