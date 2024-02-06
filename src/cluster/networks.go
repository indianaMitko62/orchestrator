package cluster

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type OrchNetwork struct {
	cli               *client.Client
	Name              string
	ID                string
	Status            string
	ActiveConnections int
	NetworkConfig     types.NetworkCreate
}

func (network *OrchNetwork) NetwCreate(opts types.NetworkCreate) (types.NetworkCreateResponse, error) {
	res, err := network.cli.NetworkCreate(context.Background(), network.Name, opts)
	if err != nil {
		slog.Error("Could not create network", "name", network.Name)
		return res, err
	}
	network.Status = "created"
	network.ID = res.ID
	return res, err
}

func (network *OrchNetwork) NetwConnect(container OrchContainer, config *network.EndpointSettings) error {
	err := network.cli.NetworkConnect(context.Background(), network.Name, container.ContID, config)
	if err != nil {
		slog.Error("Could not connect container to network", "container", container.ContainerConfig.Hostname, "network", network.Name)
		return err
	}
	network.ActiveConnections += 1
	return err
}

func (network *OrchNetwork) NetwDisconnect(container *OrchContainer, force bool) error {
	err := network.cli.NetworkDisconnect(context.Background(), network.ID, container.ContID, force)
	if err != nil {
		slog.Error("Could not disconnect container from network", "container", container.ContainerConfig.Hostname, "network", network.Name)
		return err
	}
	network.ActiveConnections -= 1
	return err
}

func (network *OrchNetwork) NetwInspect(opts types.NetworkInspectOptions) (types.NetworkResource, error) {
	res, err := network.cli.NetworkInspect(context.Background(), network.ID, opts)
	if err != nil {
		slog.Error("Could not inspect network", "network", network.Name)
		return res, err
	}
	return res, err
}

func (network *OrchNetwork) NetwList(opts types.NetworkListOptions) ([]types.NetworkResource, error) {
	res, err := network.cli.NetworkList(context.Background(), opts)
	if err != nil {
		slog.Error("Could not list networks")
		return res, err
	}
	return res, err
}

func (network *OrchNetwork) NetwRemove() error {
	err := network.cli.NetworkRemove(context.Background(), network.ID)
	if err != nil {
		slog.Error("Could not remove network", "network", network.Name)
		return err
	}
	network.Status = "removed"
	return err
}

func (network *OrchNetwork) NetwsPrune(pruneFilters filters.Args) (types.NetworksPruneReport, error) {
	res, err := network.cli.NetworksPrune(context.Background(), pruneFilters)
	if err != nil {
		slog.Error("Could not prune networks")
		return res, err
	}
	return res, err
}
