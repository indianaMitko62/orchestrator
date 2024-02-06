package cluster

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

type OrchVolume struct {
	cli    *client.Client
	Name   string
	Status string
	Config volume.CreateOptions
}

/*
TODO: mv List and Prune
*/

func (volume *OrchVolume) CreateVol(opts volume.CreateOptions) (volume.Volume, error) {
	slog.Info("Creating volume", "name", volume.Name)
	vol, err := volume.cli.VolumeCreate(context.Background(), opts)
	if err != nil {
		slog.Error("could not create volume", "name", volume.Name)
		return vol, err
	}
	volume.Status = "created"
	slog.Info("Volume created", "name", &volume.Name)
	return vol, err
}

func (volume *OrchVolume) InspectVol() (volume.Volume, error) {
	slog.Info("Inspecting volume", "name", volume.Name)
	vol, err := volume.cli.VolumeInspect(context.Background(), volume.Name)
	if err != nil {
		slog.Error("could not inspect volume", "name", volume.Name)
		return vol, err
	}
	slog.Info("Volume inspected", "name", &volume.Name)
	return vol, err
}

func (volume *OrchVolume) ListVols(opts volume.ListOptions) (volume.ListResponse, error) {
	slog.Info("Listing volumes")
	res, err := volume.cli.VolumeList(context.Background(), opts)
	if err != nil {
		slog.Error("could not list volumes")
		return res, err
	}
	slog.Info("Volumes listed")
	return res, err
}

func (volume *OrchVolume) RemoveVol(force bool) error {
	slog.Info("Removing volume", "name", volume.Name)
	err := volume.cli.VolumeRemove(context.Background(), volume.Name, force)
	if err != nil {
		slog.Error("could not list volumes")
		return err
	}
	volume.Status = "removed"
	slog.Info("Volume removed", "name", &volume.Name)
	return err
}

func (volume *OrchVolume) PruneVols(pruneFilters filters.Args) (types.VolumesPruneReport, error) {
	slog.Info("Pruning volume")
	report, err := volume.cli.VolumesPrune(context.Background(), pruneFilters)
	if err != nil {
		slog.Error("could not list volumes")
		return report, err
	}
	slog.Info("Volumes pruned", "name", &volume.Name)
	return report, err
}
