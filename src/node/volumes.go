package node

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

type OrchVolume struct {
	Name       string
	Status     string
	MountPoint string
}

func (n *NodeService) VolCreate(volume *OrchVolume, opts volume.CreateOptions) (volume.Volume, error) {
	vol, err := n.cli.VolumeCreate(context.Background(), opts)
	if err != nil {
		slog.Error("could not create volume", "name", volume.Name)
		return vol, err
	}
	volume.Status = "created"

	return vol, err
}

func (n *NodeService) VolInspect(volume *OrchVolume) (volume.Volume, error) {
	vol, err := n.cli.VolumeInspect(context.Background(), volume.Name)
	if err != nil {
		slog.Error("could not inspect volume", "name", volume.Name)
		return vol, err
	}
	volume.Status = "created"
	return vol, err
}

func (n *NodeService) VolList(opts volume.ListOptions) (volume.ListResponse, error) {
	res, err := n.cli.VolumeList(context.Background(), opts)
	if err != nil {
		slog.Error("could not list volumes")
		return res, err
	}
	return res, err
}

func (n *NodeService) VolRemove(volume *OrchVolume, force bool) error {
	err := n.cli.VolumeRemove(context.Background(), volume.Name, force)
	if err != nil {
		slog.Error("could not list volumes")
		return err
	}
	volume.Status = "removed"
	return err
}

func (n *NodeService) VolPrune(pruneFilters filters.Args) (types.VolumesPruneReport, error) {
	report, err := n.cli.VolumesPrune(context.Background(), pruneFilters)
	if err != nil {
		slog.Error("could not list volumes")
		return report, err
	}
	return report, err
}
