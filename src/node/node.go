package node

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types"
)

func (n *NodeService) ListContainers(Cont *Container, Opts types.ContainerListOptions) ([]types.Container, error) {
	slog.Info("Received list container request", "name", Cont.ContainerName)
	var containerList []types.Container
	containerList, err := n.cli.ContainerList(context.Background(), Opts)
	if err != nil {
		slog.Error("could not list containers", "name", Cont.ContainerName)
		return nil, err
	}
	return containerList, nil
}

//disk usage and others functions managing overall node configuration
