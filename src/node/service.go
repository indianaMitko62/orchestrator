package node

import (
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
			Port:             ":1986",
			LogsPath:         "/logs",
			ClusterStatePath: "/clusterState",
			NodeStatusPath:   "/nodeStatus",
		},
		DesiredNodeState: cluster.NewNodeState(),
		clusterChangeLog: cluster.NewLog("./logs/nodeLogs/clusterChangeLog"),
		nodeLog:          cluster.NewLog("./logs/nodeLogs/nodeLog"),
	}
	return ns, nil
}
