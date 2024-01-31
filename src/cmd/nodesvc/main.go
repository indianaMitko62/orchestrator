package main

import (
	"log/slog"

	"github.com/indianaMitko62/orchestrator/src/node"
	"github.com/indianaMitko62/orchestrator/src/noderpc"
)

func main() {
	var err error
	slog.Info("before rpc")
	go noderpc.Listen()
	// if err != nil {
	// 	slog.Error("could not start listener", "err", err)
	// 	os.Exit(1)
	// }
	slog.Info("after rpc")
	nsvc, err := node.NewNodeService()
	if err != nil {
		slog.Error("could not create Node service")
	}
	nsvc.Node()

	// add InfrastructureState
}
