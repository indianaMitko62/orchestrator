package node

import (
	"log/slog"
	"os"

	"github.com/docker/docker/client"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type NodeService struct {
	cluster.NodeSettings
	cli              *client.Client
	DesiredNodeState *cluster.NodeState // no pointer required
	CurrentNodeState *cluster.NodeState
	clusterChangeLog *cluster.Log
	nodeLog          *cluster.Log
	LogsDir          string
}

func NewNodeService() (*NodeService, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	ns := &NodeService{
		cli: cli,
		NodeSettings: cluster.NodeSettings{
			Name:             "Node1",
			Address:          "127.0.0.1", // Node IP from machine setup. Left to 127.0.0.1 for testing purposes.
			MasterAddress:    "127.0.0.1",
			MasterPort:       ":1986",
			LogsPath:         "/logs",
			ClusterStatePath: "/clusterState",
			NodeStatusPath:   "/nodeStatus",
		},
		DesiredNodeState: cluster.NewNodeState(),
	}
	ns.LogsDir = "./logs/" + ns.Name + "Logs/"
	if err := os.Mkdir(ns.LogsDir, 0755); os.IsExist(err) {
		slog.Info("Directory exists", "name", ns.LogsDir)
	} else {
		slog.Info("Directory created", "name", ns.LogsDir)
	}
	ns.clusterChangeLog = cluster.NewLog(ns.LogsDir + "clusterChangeLog")
	ns.nodeLog = cluster.NewLog(ns.LogsDir + "nodeLog")
	return ns, nil
}
