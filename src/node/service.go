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

func NewLog(file string) *Log {
	logFile, _ := os.Create(file)
	logReader, err := os.Open(file)
	if err != nil {
		slog.Error("Could not create logger")
	}
	logWriter := io.MultiWriter(os.Stdout, logFile)
	logger := slog.New(slog.NewTextHandler(logWriter, nil))
	return &Log{
		Logger:    logger,
		logReader: logReader,
	}
}

func NewNodeService() (*NodeService, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	ns := &NodeService{
		cli: cli,
		NodeSettings: cluster.NodeSettings{
			Name:    "Node1",
			Address: "127.0.0.1", // Node IP from machine setup. Left to 127.0.0.1 for testing purposes.
		},
		DesiredNodeState: cluster.NewNodeState(),
		clusterChangeLog: NewLog("./clusterChangeLog"),
		nodeLog:          NewLog("./nodeLog"),
	}
	return ns, nil
}
