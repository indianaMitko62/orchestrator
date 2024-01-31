package node

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
)

type OrchNetwork struct {
	Name   string
	Status string
	ID     string
}

func (n *NodeService) NetwCreate(network *OrchNetwork, opts types.NetworkCreate) (types.NetworkCreateResponse, error) {
	res, err := n.cli.NetworkCreate(context.Background(), network.Name, opts)
	if err != nil {
		slog.Error("Could not create network", "name", network.Name)
		return res, err
	}
	network.Status = "created"
	network.ID = res.ID
	return res, err
}

func (n *NodeService) NetwConnect(network *OrchNetwork, container OrchContainer, config *network.EndpointSettings) error {
	err := n.cli.NetworkConnect(context.Background(), network.Name, container.ContID, config)
	if err != nil {
		slog.Error("Could not connect container to network", "container", container.Name, "network", network.Name)
		return err
	}
	network.Status = "connected to"
	return err
}

func (n *NodeService) NetwDisconnect(network *OrchNetwork, container *OrchContainer, force bool) error {
	err := n.cli.NetworkDisconnect(context.Background(), network.ID, container.ContID, force)
	if err != nil {
		slog.Error("Could not disconnect container from network", "container", container.Name, "network", network.Name)
		return err
	}
	network.Status = "created"
	return err
}

func (n *NodeService) NetwInspect(network OrchNetwork, opts types.NetworkInspectOptions) (types.NetworkResource, error) {
	res, err := n.cli.NetworkInspect(context.Background(), network.ID, opts)
	if err != nil {
		slog.Error("Could not inspect network", "network", network.Name)
		return res, err
	}
	return res, err
}

func (n *NodeService) NetwList(opts types.NetworkListOptions) ([]types.NetworkResource, error) {
	res, err := n.cli.NetworkList(context.Background(), opts)
	if err != nil {
		slog.Error("Could not list networks")
		return res, err
	}
	return res, err
}

func (n *NodeService) NetwRemove(network *OrchNetwork) error {
	err := n.cli.NetworkRemove(context.Background(), network.ID)
	if err != nil {
		slog.Error("Could not remove network", "network", network.Name)
		return err
	}
	network.Status = "removed"
	return err
}

func (n *NodeService) NetwsPrune(pruneFilters filters.Args) (types.NetworksPruneReport, error) {
	res, err := n.cli.NetworksPrune(context.Background(), pruneFilters)
	if err != nil {
		slog.Error("Could not prune networks")
		return res, err
	}
	return res, err
}
