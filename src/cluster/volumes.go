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

func (volume *OrchVolume) VolCreate(opts volume.CreateOptions) (volume.Volume, error) {
	vol, err := volume.cli.VolumeCreate(context.Background(), opts)
	if err != nil {
		slog.Error("could not create volume", "name", volume.Name)
		return vol, err
	}
	volume.Status = "created"

	return vol, err
}

func (volume *OrchVolume) VolInspect() (volume.Volume, error) {
	vol, err := volume.cli.VolumeInspect(context.Background(), volume.Name)
	if err != nil {
		slog.Error("could not inspect volume", "name", volume.Name)
		return vol, err
	}
	volume.Status = "created"
	return vol, err
}

func (volume *OrchVolume) VolList(opts volume.ListOptions) (volume.ListResponse, error) {
	res, err := volume.cli.VolumeList(context.Background(), opts)
	if err != nil {
		slog.Error("could not list volumes")
		return res, err
	}
	return res, err
}

func (volume *OrchVolume) VolRemove(force bool) error {
	err := volume.cli.VolumeRemove(context.Background(), volume.Name, force)
	if err != nil {
		slog.Error("could not list volumes")
		return err
	}
	volume.Status = "removed"
	return err
}

func (volume *OrchVolume) VolPrune(pruneFilters filters.Args) (types.VolumesPruneReport, error) {
	report, err := volume.cli.VolumesPrune(context.Background(), pruneFilters)
	if err != nil {
		slog.Error("could not list volumes")
		return report, err
	}
	return report, err
}
