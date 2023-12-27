package main

import (
	"log/slog"
	"os"

	"github.com/indianaMitko62/orchestrator/src/master"
)

func main() {
	msvc := master.NewMasterService([]*master.NodeSettings{
		{
			Name:    "local node",
			Address: "localhost:1234",
		},
	})

	err := msvc.ConnectToNodes()
	if err != nil {
		slog.Error("could not connect to nodes", "err", err)
		os.Exit(1)
	}

	err = msvc.CreateCont()
	if err != nil {
		slog.Error("could not say hello to world", "err", err)
		os.Exit(1)
	}
}
