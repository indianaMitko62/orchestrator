package node

import (
	"io"
	"log/slog"
	"os"

	"github.com/docker/docker/client"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type NodeService struct {
	cluster.NodeSettings
	cli              *client.Client
	DesiredNodeState *cluster.NodeState
	CurrentNodeState *cluster.NodeState
	clusterChangeLog *slog.Logger
	logReader        io.Reader
}

func NewNodeService() (*NodeService, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	ns := &NodeService{
		cli: cli,
		NodeSettings: cluster.NodeSettings{ // Can have separate init function for visability
			Name:    "Node1",
			Address: "127.0.0.1", // Node IP from machine setup. Left to 127.0.0.1 for testing purposes.
		},
		DesiredNodeState: cluster.NewNodeState(),
	}
	f, _ := os.Create("./dat2")
	ns.logReader, err = os.Open("./dat2")
	if err != nil {
		slog.Error("Could not create logger")
	}
	logWriter := io.MultiWriter(os.Stdout, f)
	ns.clusterChangeLog = slog.New(slog.NewTextHandler(logWriter, nil))
	return ns, nil
}
