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
	ContID           string
	ContStatus       string
	ContainerConfig  *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
}

/*
NOTES: Most of these functionalities do not need to be accessable remotely.
*/

// type ContainerSettings struct { // not needed at the moment, but who knows
// 	Cont *Container
// }

func (cont *OrchContainer) CreateCont() (string, error) {
	slog.Info("Received create request", "name", cont.ContainerConfig.Hostname)
	reply, err := cont.cli.ContainerCreate(context.Background(),
		cont.ContainerConfig,
		cont.HostConfig,
		cont.NetworkingConfig,
		nil,
		cont.ContainerConfig.Hostname)
	if err != nil {
		slog.Error("could not create container", "name", cont.ContainerConfig.Hostname)
		return "", err
	}

	cont.ContID = reply.ID
	cont.ContStatus = "created"
	slog.Info("Container created", "name", &cont.ContainerConfig.Hostname)
	return reply.ID, nil
}

func (cont *OrchContainer) StartCont(Opts types.ContainerStartOptions) error {
	slog.Info("Received start request", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerStart(context.Background(), cont.ContID, Opts)
	if err != nil {
		slog.Error("could not start container", "name", cont.ContID)
		return err
	}
	cont.ContStatus = "running"
	slog.Info("Container started", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (cont *OrchContainer) StopCont(Opts container.StopOptions) error {
	slog.Info("Received stop request", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerStop(context.Background(), cont.ContID, Opts)
	if err != nil {
		slog.Error("could not stop container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.ContStatus = "stopped"
	slog.Info("Container stopped", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (cont *OrchContainer) InspectCont(getSize bool) (types.ContainerJSON, []byte, error) {
	slog.Info("Received inspect request", "name", cont.ContainerConfig.Hostname)
	json, byte, err := cont.cli.ContainerInspectWithRaw(context.Background(), cont.ContID, getSize)
	if err != nil {
		slog.Error("could not start container", "name", cont.ContID)
		return json, byte, err
	}
	slog.Info("Container inspected", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return json, byte, nil
}

func (cont *OrchContainer) LogCont(Opts types.ContainerLogsOptions) (io.ReadCloser, error) {
	slog.Info("Received log request", "name", cont.ContainerConfig.Hostname)
	out, err := cont.cli.ContainerLogs(context.Background(), cont.ContID, Opts)
	if err != nil {
		slog.Error("could not log container", "name", cont.ContainerConfig.Hostname)
		return nil, err
	}
	slog.Info("Containers Logs returned", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return out, nil
}

func (cont *OrchContainer) KillCont(signal string) error {
	slog.Info("Received kill request", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerKill(context.Background(), cont.ContID, signal)
	if err != nil {
		slog.Error("could not kill container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.ContStatus = "killed"
	slog.Info("Container killed", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (cont *OrchContainer) RemoveCont(Opts types.ContainerRemoveOptions) error {
	slog.Info("Received remove request", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerRemove(context.Background(), cont.ContID, Opts)
	if err != nil {
		slog.Error("could not remove container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.ContStatus = "removed"
	slog.Info("Container removed", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (cont *OrchContainer) ListContainers(Opts types.ContainerListOptions) ([]types.Container, error) {
	slog.Info("Received list container request", "name", cont.ContainerConfig.Hostname)
	var containerList []types.Container
	containerList, err := cont.cli.ContainerList(context.Background(), Opts)
	if err != nil {
		slog.Error("could not list containers", "name", cont.ContainerConfig.Hostname)
		return nil, err
	}
	return containerList, nil
}

func (cont *OrchContainer) PauseCont() error {
	slog.Info("Received pause request", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerPause(context.Background(), cont.ContID)
	if err != nil {
		slog.Error("could not pause container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.ContStatus = "paused"
	slog.Info("Container paused", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (cont *OrchContainer) UnpauseCont() error {
	slog.Info("Received unpause request", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.ContainerUnpause(context.Background(), cont.ContID)
	if err != nil {
		slog.Error("could not unpause container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.ContStatus = "running"
	slog.Info("Container unpaused", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (cont *OrchContainer) CopyToCont(dest string, src io.Reader, Opts types.CopyToContainerOptions) error {
	slog.Info("Received copyTo request", "name", cont.ContainerConfig.Hostname)
	err := cont.cli.CopyToContainer(context.Background(), cont.ContID, dest, src, Opts)
	if err != nil {
		slog.Error("could not copyTo container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	slog.Info("copyTo Container", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (cont *OrchContainer) CopyFromCont(src string) (io.ReadCloser, error) {
	slog.Info("Received copyFrom request", "name", cont.ContainerConfig.Hostname)
	res, _, err := cont.cli.CopyFromContainer(context.Background(), cont.ContID, src)
	if err != nil {
		slog.Error("could not copyFrom container", "name", cont.ContainerConfig.Hostname)
		return nil, err
	}
	slog.Info("copyFrom Container", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return res, nil
}

func (cont *OrchContainer) TopCont(args []string) (container.ContainerTopOKBody, error) {
	slog.Info("Received Top request", "name", cont.ContainerConfig.Hostname)
	res, err := cont.cli.ContainerTop(context.Background(), cont.ContID, args)
	if err != nil {
		slog.Error("could not top container", "name", cont.ContainerConfig.Hostname)
		return res, err
	}
	slog.Info("container top returned", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return res, nil
}

func (cont *OrchContainer) StatCont(stream bool) (types.ContainerStats, error) {
	slog.Info("Received stat request", "name", cont.ContainerConfig.Hostname)
	res, err := cont.cli.ContainerStats(context.Background(), cont.ContID, stream)
	if err != nil {
		slog.Error("could not stat container", "name", cont.ContainerConfig.Hostname)
		return res, err
	}
	slog.Info("container stat returned", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return res, nil
}
