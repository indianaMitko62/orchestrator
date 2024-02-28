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
	Cli               *client.Client
	Name              string
	ID                string
	CurrentStatus     string
	DesiredStatus     string
	ActiveConnections int
	NetworkConfig     types.NetworkCreate
}

/*
TODO: mv List and Prune ??
*/

func (network *OrchNetwork) CreateNet(opts types.NetworkCreate) (string, error) {
	slog.Info("Creating network", "name", network.Name)
	res, err := network.Cli.NetworkCreate(context.Background(), network.Name, opts)
	if err != nil {
		slog.Error("Could not create network", "name", network.Name, "error", err)
		return "", err
	}
	network.CurrentStatus = "created"
	network.ID = res.ID
	slog.Info("Network Created", "name", network.Name, "ID", network.ID)
	return res.ID, err
}

func (network *OrchNetwork) ConnectToNet(container OrchContainer, config *network.EndpointSettings) error {
	slog.Info("Connecting to network", "name", network.Name)
	err := network.Cli.NetworkConnect(context.Background(), network.Name, container.ContainerConfig.Hostname, config)
	if err != nil {
		slog.Error("Could not connect container to network", "container", container.ContainerConfig.Hostname, "network", network.Name, "err", err.Error())
		return err
	}
	network.ActiveConnections += 1
	slog.Info("Connected to network", "name", network.Name, "ID", network.ID, "container", container.ContainerConfig.Hostname)
	return err
}

func (network *OrchNetwork) DisconnectFromNet(container OrchContainer, force bool) error {
	slog.Info("Disconnecting from network", "name", network.Name)
	err := network.Cli.NetworkDisconnect(context.Background(), network.Name, container.ContainerConfig.Hostname, force)
	if err != nil {
		slog.Error("Could not disconnect container from network", "container", container.ContainerConfig.Hostname, "network", network.Name)
		return err
	}
	network.ActiveConnections -= 1
	slog.Info("Disconnected from network", "name", network.Name, "ID", network.ID, "container", container.ContainerConfig.Hostname)
	return err
}

func (network *OrchNetwork) InspectNet(opts types.NetworkInspectOptions) (types.NetworkResource, error) {
	slog.Info("Inspecting networks", "name", network.Name)
	res, err := network.Cli.NetworkInspect(context.Background(), network.Name, opts)
	if err != nil {
		slog.Error("Could not inspect network", "network", network.Name)
		return res, err
	}
	slog.Info("Network inspected", "name", network.Name, "ID", network.ID)
	return res, err
}

func (network *OrchNetwork) ListNets(opts types.NetworkListOptions) ([]types.NetworkResource, error) {
	slog.Info("Listing networks")
	res, err := network.Cli.NetworkList(context.Background(), opts)
	if err != nil {
		slog.Error("Could not list networks")
		return res, err
	}
	slog.Info("Networks listed")
	return res, err
}

func (network *OrchNetwork) RemoveNet() error {
	slog.Info("Removing network", "name", network.Name)
	err := network.Cli.NetworkRemove(context.Background(), network.Name)
	if err != nil {
		slog.Error("Could not remove network", "network", network.Name, "err", err.Error())
		return err
	}
	network.CurrentStatus = "removed"
	slog.Info("Network removed", "name", network.Name, "ID", network.ID)
	return err
}

func (network *OrchNetwork) PruneNetws(pruneFilters filters.Args) (types.NetworksPruneReport, error) {
	slog.Info("Pruning networks")
	res, err := network.Cli.NetworksPrune(context.Background(), pruneFilters)
	if err != nil {
		slog.Error("Could not prune networks")
		return res, err
	}
	slog.Info("Networks pruned")
	return res, err
}
