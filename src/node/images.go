package node

import (
	"context"
	"io"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
)

type OrchImage struct {
	Name    *string
	Tag     *string
	Version float32
	Status  *string
}

/*
TODO: All of the image functionalities are to be moved to Master, except ImgPull, ImgList, ImgTag(probably), ImgRemove, ImgInspect.
NOTES: For I do not believe they have to be accessable remotely.
*/

func (n *NodeService) ImgBuild(image *OrchImage, buildContext io.Reader, opts types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	res, err := n.cli.ImageBuild(context.Background(), buildContext, opts)
	if err != nil {
		slog.Error("could not build image", "name", image.Name)
		return res, err
	}
	*image.Status = "built"
	*image.Name = opts.Tags[0]
	return res, nil
} // ImageCreate???

func (n *NodeService) ImgPull(name string, opts types.ImagePullOptions) (io.ReadCloser, error) {
	res, err := n.cli.ImagePull(context.Background(), name, opts)
	if err != nil {
		slog.Error("could not pull image", "name", name)
		return res, err
	}
	return res, nil
}

func (n *NodeService) ImgPush(image *OrchImage, opts types.ImagePushOptions) (io.ReadCloser, error) {
	res, err := n.cli.ImagePush(context.Background(), *image.Name, opts)
	if err != nil {
		slog.Error("could not push image", "name", image.Name)
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

func (n *NodeService) ImgHist(image *OrchImage) ([]image.HistoryResponseItem, error) {
	res, err := n.cli.ImageHistory(context.Background(), *image.Name)
	if err != nil {
		slog.Error("could not get history for image", "name", image.Name)
		return res, err
	}
	return res, nil
}

func (n *NodeService) ImgSave(image *OrchImage) (io.ReadCloser, error) {
	ids := []string{*image.Name}
	res, err := n.cli.ImageSave(context.Background(), ids)
	if err != nil {
		slog.Error("could not save image", "name", image.Name)
		return res, err
	}
	return res, nil
}

func (n *NodeService) ImgLoad(image *OrchImage, input io.ReadCloser, quiet bool) (types.ImageLoadResponse, error) {
	res, err := n.cli.ImageLoad(context.Background(), input, quiet) // quiet for minimal output (just the image id)
	if err != nil {
		slog.Error("could not load image", "name", image.Name)
		return res, err
	}
	return res, nil
}

func (n *NodeService) ImgRemove(image *OrchImage, opts types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	res, err := n.cli.ImageRemove(context.Background(), *image.Name, opts)
	if err != nil {
		slog.Error("could not remove image", "name", image.Name)
		return res, err
	}
	*image.Status = "removed"
	return res, nil
}

func (n *NodeService) ImgInspect(image OrchImage) (types.ImageInspect, []byte, error) {
	res, raw, err := n.cli.ImageInspectWithRaw(context.Background(), *image.Name)
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
