package main

import (
	"log/slog"
	"os"

	"github.com/indianaMitko62/orchestrator/src/cluster"
	"github.com/indianaMitko62/orchestrator/src/node"
	"gopkg.in/yaml.v3"
)

func main() {
	var err error
	confFile := os.Args[1]
	if confFile == "" {
		slog.Error("No command line argument")
	}
	f, err := os.Open(confFile)
	if err != nil {
		slog.Error("Could not open config file", "name", confFile)
	}

	var nodeSet cluster.NodeSettings
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&nodeSet)
	if err != nil {
		slog.Error("Could not decode config file", "name", confFile)
	}
	f.Close()
	nsvc, err := node.NewNodeService(nodeSet)
	if err != nil {
		slog.Error("could not create Node service")
	}
	nsvc.Node()
}
