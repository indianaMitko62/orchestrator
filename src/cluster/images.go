package cluster

import (
	"context"
	"io"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type OrchImage struct {
	cli          *client.Client
	Name         string
	Tag          string
	ID           string
	Status       string
	BuildOptions types.ImageBuildOptions
	// ImagePull options. Cannot marshall RequestPrivilegeFunc in ImagePull
	All          bool
	RegistryAuth string
	Platform     string
}

/*
NOTES: For I do not believe they have to be accessable remotely.
*/

func (img *OrchImage) ImgBuild(buildContext io.Reader, opts types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	res, err := img.cli.ImageBuild(context.Background(), buildContext, opts)
	if err != nil {
		slog.Error("could not build image", "name", img.Name)
		return res, err
	}
	img.Status = "built"
	img.Name = opts.Tags[0]
	return res, nil
} // ImageCreate???

func (img *OrchImage) ImgPull(name string, opts types.ImagePullOptions) (io.ReadCloser, error) {
	res, err := img.cli.ImagePull(context.Background(), name, opts)
	if err != nil {
		slog.Error("could not pull image", "name", name)
		return res, err
	}
	return res, nil
}

func (img *OrchImage) ImgPush(opts types.ImagePushOptions) (io.ReadCloser, error) {
	res, err := img.cli.ImagePush(context.Background(), img.Name, opts)
	if err != nil {
		slog.Error("could not push image", "name", img.Name)
		return res, err
	}
	return res, nil
}

func (img *OrchImage) ImgList(opts types.ImageListOptions) ([]types.ImageSummary, error) {
	res, err := img.cli.ImageList(context.Background(), opts)
	if err != nil {
		slog.Error("could not list images")
		return res, err
	}
	return res, nil
}

func (img *OrchImage) ImgTag(src string, target string) error {
	err := img.cli.ImageTag(context.Background(), src, target)
	if err != nil {
		slog.Error("could not tag image", "name", img.Name)
		return err
	}
	return nil
}

func (img *OrchImage) ImgHist() ([]image.HistoryResponseItem, error) {
	res, err := img.cli.ImageHistory(context.Background(), img.Name)
	if err != nil {
		slog.Error("could not get history for image", "name", img.Name)
		return res, err
	}
	return res, nil
}

func (img *OrchImage) ImgSave() (io.ReadCloser, error) {
	ids := []string{img.Name}
	res, err := img.cli.ImageSave(context.Background(), ids)
	if err != nil {
		slog.Error("could not save image", "name", img.Name)
		return res, err
	}
	return res, nil
}

func (img *OrchImage) ImgLoad(input io.ReadCloser, quiet bool) (types.ImageLoadResponse, error) {
	res, err := img.cli.ImageLoad(context.Background(), input, quiet) // quiet for minimal output (just the image id)
	if err != nil {
		slog.Error("could not load image", "name", img.Name)
		return res, err
	}
	return res, nil
}

func (img *OrchImage) ImgRemove(opts types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	res, err := img.cli.ImageRemove(context.Background(), img.Name, opts)
	if err != nil {
		slog.Error("could not remove image", "name", img.Name)
		return res, err
	}
	img.Status = "removed"
	return res, nil
}

func (img *OrchImage) ImgInspect() (types.ImageInspect, []byte, error) {
	res, raw, err := img.cli.ImageInspectWithRaw(context.Background(), img.Name)
	if err != nil {
		slog.Error("could not inspect image", "name", img.Name)
		return res, raw, err
	}
	return res, raw, nil
}

func (img *OrchImage) ImgPrune(pruneFilters filters.Args) (types.ImagesPruneReport, error) {
	res, err := img.cli.ImagesPrune(context.Background(), pruneFilters)
	if err != nil {
		slog.Error("could not prune images")
		return res, err
	}
	return res, nil
}
