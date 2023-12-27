package main

import (
	"log/slog"
	"os"

	"github.com/indianaMitko62/orchestrator/src/orchestrator/noderpc"
)

func main() {
	err := noderpc.Listen()
	if err != nil {
		slog.Error("could not start listener", "err", err)
		os.Exit(1)
	}
}
