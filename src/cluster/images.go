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
TODO: mv List and Prune
NOTES: For I do not believe they have to be accessable remotely.
*/

func (img *OrchImage) BuildImg(buildContext io.Reader, opts types.ImageBuildOptions) (types.ImageBuildResponse, error) {
	slog.Info("Building image", "name", img.Name)
	res, err := img.cli.ImageBuild(context.Background(), buildContext, opts)
	if err != nil {
		slog.Error("could not build image", "name", img.Name)
		return res, err
	}
	img.Status = "built"
	img.Name = opts.Tags[0]
	slog.Info("Image built", "name", img.Name, "ID", img.ID)
	return res, nil
} // ImageCreate???

func (img *OrchImage) PullImg(opts types.ImagePullOptions) (io.ReadCloser, error) {
	slog.Info("Pulling image", "name", img.Name)
	res, err := img.cli.ImagePull(context.Background(), img.Name, opts)
	if err != nil {
		slog.Error("could not pull image", "name", img.Name)
		return res, err
	}
	slog.Info("Image pulled", "name", img.Name, "ID", img.ID)
	return res, nil
}

func (img *OrchImage) PushImg(opts types.ImagePushOptions) (io.ReadCloser, error) {
	slog.Info("Pushing image", "name", img.Name)
	res, err := img.cli.ImagePush(context.Background(), img.Name, opts)
	if err != nil {
		slog.Error("could not push image", "name", img.Name, "ID", img.ID)
		return res, err
	}
	slog.Info("Image pushed", "name", img.Name)
	return res, nil
}

func (img *OrchImage) ListImg(opts types.ImageListOptions) ([]types.ImageSummary, error) {
	slog.Info("Listing images")
	res, err := img.cli.ImageList(context.Background(), opts)
	if err != nil {
		slog.Error("could not list images")
		return res, err
	}
	slog.Info("Images listed")
	return res, nil
}

func (img *OrchImage) TagImg(src string, target string) error {
	slog.Info("Tagging image", "name", img.Name)
	err := img.cli.ImageTag(context.Background(), src, target)
	if err != nil {
		slog.Error("could not tag image", "name", img.Name)
		return err
	}
	slog.Info("Image tagged", "name", img.Name, "ID", img.ID)
	return nil
}

func (img *OrchImage) HistImg() ([]image.HistoryResponseItem, error) {
	slog.Info("Getting image history", "name", img.Name)
	res, err := img.cli.ImageHistory(context.Background(), img.Name)
	if err != nil {
		slog.Error("could not get history for image", "name", img.Name)
		return res, err
	}
	slog.Info("Got image history", "name", img.Name, "ID", img.ID)
	return res, nil
}

func (img *OrchImage) SaveImg() (io.ReadCloser, error) {
	slog.Info("Saving image", "name", img.Name)
	ids := []string{img.Name}
	res, err := img.cli.ImageSave(context.Background(), ids)
	if err != nil {
		slog.Error("could not save image", "name", img.Name)
		return res, err
	}
	slog.Info("Image saved", "name", img.Name, "ID", img.ID)
	return res, nil
}

func (img *OrchImage) LoadImg(input io.ReadCloser, quiet bool) (types.ImageLoadResponse, error) {
	slog.Info("Loading image", "name", img.Name)
	res, err := img.cli.ImageLoad(context.Background(), input, quiet) // quiet for minimal output (just the image id)
	if err != nil {
		slog.Error("could not load image", "name", img.Name)
		return res, err
	}
	slog.Info("Image loaded", "name", img.Name, "ID", img.ID)
	return res, nil
}

func (img *OrchImage) RemoveImg(opts types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error) {
	slog.Info("Removing image", "name", img.Name)
	res, err := img.cli.ImageRemove(context.Background(), img.Name, opts)
	if err != nil {
		slog.Error("could not remove image", "name", img.Name)
		return res, err
	}
	img.Status = "removed"
	slog.Info("Image removed", "name", img.Name, "ID", img.ID)
	return res, nil
}

func (img *OrchImage) InspectImg() (types.ImageInspect, []byte, error) {
	slog.Info("Inspecting image", "name", img.Name)
	res, raw, err := img.cli.ImageInspectWithRaw(context.Background(), img.Name)
	if err != nil {
		slog.Error("could not inspect image", "name", img.Name)
		return res, raw, err
	}
	slog.Info("Image inspected", "name", img.Name, "ID", img.ID)
	return res, raw, nil
}

func (img *OrchImage) PruneImgs(pruneFilters filters.Args) (types.ImagesPruneReport, error) {
	slog.Info("Pruning images")
	res, err := img.cli.ImagesPrune(context.Background(), pruneFilters)
	if err != nil {
		slog.Error("could not prune images")
		return res, err
	}
	slog.Info("Images pruned")
	return res, nil
}
