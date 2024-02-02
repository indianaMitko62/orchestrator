package node

import (
	"context"
	"io"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

type OrchImage struct {
	Name    string
	Tag     string
	Version float32
	Status  string
}

/*
NOTES: For I do not believe they have to be accessable remotely.
*/

func (n *NodeService) ImgPull(name string, opts types.ImagePullOptions) (io.ReadCloser, error) {
	res, err := n.cli.ImagePull(context.Background(), name, opts)
	if err != nil {
		slog.Error("could not pull image", "name", name)
		return res, err
	}
	return res, nil
}

func (n *NodeService) ImgList(opts types.ImageListOptions) ([]types.ImageSummary, error) {
	res, err := n.cli.ImageList(context.Background(), opts)
	if err != nil {
		slog.Error("could not list images")
		return res, err
	}
	return res, nil
}

func (n *NodeService) ImgTag(image *OrchImage, src string, target string) error {
	err := n.cli.ImageTag(context.Background(), src, target)
	if err != nil {
		slog.Error("could not tag image", "name", image.Name)
		return err
	}
	return nil
}

func (n *NodeService) ImgRemove(image *OrchImage, opts types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	res, err := n.cli.ImageRemove(context.Background(), image.Name, opts)
	if err != nil {
		slog.Error("could not remove image", "name", image.Name)
		return res, err
	}
	image.Status = "removed"
	return res, nil
}

func (n *NodeService) ImgInspect(image OrchImage) (types.ImageInspect, []byte, error) {
	res, raw, err := n.cli.ImageInspectWithRaw(context.Background(), image.Name)
	if err != nil {
		slog.Error("could not inspect image", "name", image.Name)
		return res, raw, err
	}
	return res, raw, nil
}

func (n *NodeService) ImgPrune(pruneFilters filters.Args) (types.ImagesPruneReport, error) {
	res, err := n.cli.ImagesPrune(context.Background(), pruneFilters)
	if err != nil {
		slog.Error("could not prune images")
		return res, err
	}
	return res, nil
}
