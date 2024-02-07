package node

import (
	"github.com/docker/docker/client"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type NodeService struct {
	cluster.NodeSettings
	cli                  *client.Client
	DesiredNodeState     *cluster.NodeState
	CurrentNodeState     *cluster.NodeState
	ClusterChangeOutcome *ClusterChangeOutcome
}

func NewNodeService() (*NodeService, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}

	return &NodeService{
		cli: cli,
		NodeSettings: cluster.NodeSettings{
			Name:    "Node1",
			Address: "127.0.0.1", // Node IP from machine setup. Left to 127.0.0.1 for testing purposes
		},
		DesiredNodeState: cluster.NewNodeState(),
	}, nil
}
