package node

import (
	"io"
	"log/slog"
	"os"

	"github.com/docker/docker/client"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type Log struct {
	Logger    *slog.Logger
	logReader io.Reader
}

type NodeService struct {
	cluster.NodeSettings
	cli              *client.Client
	DesiredNodeState *cluster.NodeState
	CurrentNodeState *cluster.NodeState
	clusterChangeLog *Log
	nodeLog          *Log
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
		clusterChangeLog: &Log{},
		nodeLog:          &Log{},
	}
	clusterLogFile, _ := os.Create("./clusterChangeLog") // separate Log init function
	ns.clusterChangeLog.logReader, err = os.Open("./clusterChangeLog")
	if err != nil {
		slog.Error("Could not create logger")
	}
	clusterLogWriter := io.MultiWriter(os.Stdout, clusterLogFile)
	ns.clusterChangeLog.Logger = slog.New(slog.NewTextHandler(clusterLogWriter, nil))

	nodeLogFile, _ := os.Create("./nodeLog")
	ns.nodeLog.logReader, err = os.Open("./nodeLog")
	if err != nil {
		slog.Error("Could not create logger")
	}
	nodeLogWriter := io.MultiWriter(os.Stdout, nodeLogFile)
	ns.nodeLog.Logger = slog.New(slog.NewTextHandler(nodeLogWriter, nil))

	return ns, nil
}
