package main

import (
	"log/slog"

	"github.com/indianaMitko62/orchestrator/src/node"
)

func main() {
	var err error
	nsvc, err := node.NewNodeService()
	if err != nil {
		slog.Error("could not create Node service")
	}
	nsvc.Node()
}
