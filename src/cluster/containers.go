package cluster

import (
	"context"
	"io"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type OrchContainer struct {
	cli              *client.Client
	ID               string
	Status           string
	ContainerConfig  *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
}

/*
TODO: mv List
NOTES: Most of these functionalities do not need to be accessable remotely.
*/

// type ContainerSettings struct { // not needed at the moment, but who knows
// 	Cont *Container
// }

func (cont *OrchContainer) CreateCont() (string, error) {
	slog.Info("Creating container", "name", cont.ContainerConfig.Hostname)
	reply, err := cont.cli.ContainerCreate(context.Background(),
		cont.ContainerConfig,
		cont.HostConfig,
		cont.NetworkingConfig,
		nil,
		cont.ContainerConfig.Hostname)
	if err != nil {
		slog.Error("could not create container", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
		return "", err
	}

	cont.ID = reply.ID
	cont.Status = "created"
	slog.Info("Container created", "name", &cont.ContainerConfig.Hostname)
	return reply.ID, nil
}

func (cont *OrchContainer) StartCont(Opts types.ContainerStartOptions) error {
	slog.Info("Starting container", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerStart(context.Background(), cont.ID, Opts)
	if err != nil {
		slog.Error("could not start container", "name", cont.ID)
		return err
	}
	cont.Status = "running"
	slog.Info("Container started", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return nil
}

func (cont *OrchContainer) StopCont(Opts container.StopOptions) error {
	slog.Info("Stopping container", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerStop(context.Background(), cont.ID, Opts)
	if err != nil {
		slog.Error("could not stop container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.Status = "stopped"
	slog.Info("Container stopped", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return nil
}

func (cont *OrchContainer) InspectCont(getSize bool) (types.ContainerJSON, []byte, error) {
	slog.Info("Inspecting container", "name", cont.ContainerConfig.Hostname)
	json, byte, err := cont.cli.ContainerInspectWithRaw(context.Background(), cont.ID, getSize)
	if err != nil {
		slog.Error("could not start container", "name", cont.ID)
		return json, byte, err
	}
	slog.Info("Container inspected", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return json, byte, nil
}

func (cont *OrchContainer) LogCont(Opts types.ContainerLogsOptions) (io.ReadCloser, error) {
	slog.Info("Logging container", "name", cont.ContainerConfig.Hostname)
	out, err := cont.cli.ContainerLogs(context.Background(), cont.ID, Opts)
	if err != nil {
		slog.Error("could not log container", "name", cont.ContainerConfig.Hostname)
		return nil, err
	}
	slog.Info("Containers Logs returned", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return out, nil
}

func (cont *OrchContainer) KillCont(signal string) error {
	slog.Info("Killing container", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerKill(context.Background(), cont.ID, signal)
	if err != nil {
		slog.Error("could not kill container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.Status = "killed"
	slog.Info("Container killed", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return nil
}

func (cont *OrchContainer) RemoveCont(Opts types.ContainerRemoveOptions) error {
	slog.Info("Removing container", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerRemove(context.Background(), cont.ID, Opts)
	if err != nil {
		slog.Error("could not remove container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.Status = "removed"
	slog.Info("Container removed", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return nil
}

func (cont *OrchContainer) ListContainers(Opts types.ContainerListOptions) ([]types.Container, error) {
	slog.Info("Listing containers")
	var containerList []types.Container
	containerList, err := cont.cli.ContainerList(context.Background(), Opts)
	if err != nil {
		slog.Error("could not list containers", "name", cont.ContainerConfig.Hostname)
		return nil, err
	}
	slog.Info("Containers listed")
	return containerList, nil
}

func (cont *OrchContainer) PauseCont() error {
	slog.Info("Pausing containers", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerPause(context.Background(), cont.ID)
	if err != nil {
		slog.Error("could not pause container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.Status = "paused"
	slog.Info("Container paused", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return nil
}

func (cont *OrchContainer) UnpauseCont() error {
	slog.Info("Unpausing containers", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerUnpause(context.Background(), cont.ID)
	if err != nil {
		slog.Error("could not unpause container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.Status = "running"
	slog.Info("Container unpaused", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return nil
}

func (cont *OrchContainer) CopyToCont(dest string, src io.Reader, Opts types.CopyToContainerOptions) error {
	slog.Info("Copying to container", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.CopyToContainer(context.Background(), cont.ID, dest, src, Opts)
	if err != nil {
		slog.Error("could not copyTo container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	slog.Info("copyTo Container", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return nil
}

func (cont *OrchContainer) CopyFromCont(src string) (io.ReadCloser, error) {
	slog.Info("Copying from container", "name", cont.ContainerConfig.Hostname)
	res, _, err := cont.cli.CopyFromContainer(context.Background(), cont.ID, src)
	if err != nil {
		slog.Error("could not copyFrom container", "name", cont.ContainerConfig.Hostname)
		return nil, err
	}
	slog.Info("copyFrom Container", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return res, nil
}

func (cont *OrchContainer) TopCont(args []string) (container.ContainerTopOKBody, error) {
	slog.Info("Top in container", "name", cont.ContainerConfig.Hostname)
	res, err := cont.cli.ContainerTop(context.Background(), cont.ID, args)
	if err != nil {
		slog.Error("could not top container", "name", cont.ContainerConfig.Hostname)
		return res, err
	}
	slog.Info("container top returned", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return res, nil
}

func (cont *OrchContainer) StatCont(stream bool) (types.ContainerStats, error) {
	slog.Info("Container stat", "name", cont.ContainerConfig.Hostname)
	res, err := cont.cli.ContainerStats(context.Background(), cont.ID, stream)
	if err != nil {
		slog.Error("could not stat container", "name", cont.ContainerConfig.Hostname)
		return res, err
	}
	slog.Info("container stat returned", "name", cont.ContainerConfig.Hostname, "ID", cont.ID)
	return res, nil
}
