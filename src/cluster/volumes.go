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
	Cli           *client.Client
	Name          string
	CurrentStatus string
	DesiredStatus string
	Config        volume.CreateOptions
}

/*
TODO: mv List and Prune
*/

func (volume *OrchVolume) CreateVol(opts volume.CreateOptions) (volume.Volume, error) {
	slog.Info("Creating volume", "name", volume.Name)
	vol, err := volume.Cli.VolumeCreate(context.Background(), opts)
	if err != nil {
		slog.Error("could not create volume", "name", volume.Name)
		return vol, err
	}
	volume.CurrentStatus = "created"
	slog.Info("Volume created", "name", &volume.Name)
	return vol, err
}

func (volume *OrchVolume) InspectVol() (volume.Volume, error) {
	slog.Info("Inspecting volume", "name", volume.Name)
	vol, err := volume.Cli.VolumeInspect(context.Background(), volume.Name)
	if err != nil {
		slog.Error("could not inspect volume", "name", volume.Name)
		return vol, err
	}
	slog.Info("Volume inspected", "name", &volume.Name)
	return vol, err
}

func (volume *OrchVolume) ListVols(opts volume.ListOptions) (volume.ListResponse, error) {
	slog.Info("Listing volumes")
	res, err := volume.Cli.VolumeList(context.Background(), opts)
	if err != nil {
		slog.Error("could not list volumes")
		return res, err
	}
	slog.Info("Volumes listed")
	return res, err
}

func (volume *OrchVolume) RemoveVol(force bool) error {
	slog.Info("Removing volume", "name", volume.Name)
	err := volume.Cli.VolumeRemove(context.Background(), volume.Name, force)
	if err != nil {
		slog.Error("could not list volumes")
		return err
	}
	volume.CurrentStatus = "removed"
	slog.Info("Volume removed", "name", &volume.Name)
	return err
}

func (volume *OrchVolume) PruneVols(pruneFilters filters.Args) (types.VolumesPruneReport, error) {
	slog.Info("Pruning volume")
	report, err := volume.Cli.VolumesPrune(context.Background(), pruneFilters)
	if err != nil {
		slog.Error("could not list volumes")
		return report, err
	}
	slog.Info("Volumes pruned", "name", &volume.Name)
	return report, err
}
