package main

import (
	"log/slog"
	"os"

	"github.com/indianaMitko62/orchestrator/src/master"
)

func main() {
	var err error
	//name := "local node"
	var ClusterState master.ClusterState
	msvc := master.NewMasterService(&ClusterState)

	err = msvc.ConnectToNodes()
	if err != nil {
		slog.Error("could not connect to nodes", "err", err)
		os.Exit(1)
	}

}
