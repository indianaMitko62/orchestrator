package node

import (
	"context"
	"io"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type OrchContainer struct {
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

func (n *NodeService) CreateCont(cont *OrchContainer) (string, error) {
	slog.Info("Received create request", "name", cont.ContainerConfig.Hostname)
	reply, err := n.cli.ContainerCreate(context.Background(),
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

func (n *NodeService) StartCont(cont *OrchContainer, Opts types.ContainerStartOptions) error {
	slog.Info("Received start request", "name", cont.ContainerConfig.Hostname)
	err := n.cli.ContainerStart(context.Background(), cont.ContID, Opts)
	if err != nil {
		slog.Error("could not start container", "name", cont.ContID)
		return err
	}
	cont.ContStatus = "running"
	slog.Info("Container started", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (n *NodeService) StopCont(cont *OrchContainer, Opts container.StopOptions) error {
	slog.Info("Received stop request", "name", cont.ContainerConfig.Hostname)
	err := n.cli.ContainerStop(context.Background(), cont.ContID, Opts)
	if err != nil {
		slog.Error("could not stop container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.ContStatus = "stopped"
	slog.Info("Container stopped", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (n *NodeService) InspectCont(cont *OrchContainer, getSize bool) (types.ContainerJSON, []byte, error) {
	slog.Info("Received inspect request", "name", cont.ContainerConfig.Hostname)
	json, byte, err := n.cli.ContainerInspectWithRaw(context.Background(), cont.ContID, getSize)
	if err != nil {
		slog.Error("could not start container", "name", cont.ContID)
		return json, byte, err
	}
	slog.Info("Container inspected", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return json, byte, nil
}

func (n *NodeService) LogCont(cont *OrchContainer, Opts types.ContainerLogsOptions) (io.ReadCloser, error) {
	slog.Info("Received log request", "name", cont.ContainerConfig.Hostname)
	out, err := n.cli.ContainerLogs(context.Background(), cont.ContID, Opts)
	if err != nil {
		slog.Error("could not log container", "name", cont.ContainerConfig.Hostname)
		return nil, err
	}
	slog.Info("Containers Logs returned", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return out, nil
}

func (n *NodeService) KillCont(cont *OrchContainer, signal string) error {
	slog.Info("Received kill request", "name", cont.ContainerConfig.Hostname)
	err := n.cli.ContainerKill(context.Background(), cont.ContID, signal)
	if err != nil {
		slog.Error("could not kill container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.ContStatus = "killed"
	slog.Info("Container killed", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (n *NodeService) RemoveCont(cont *OrchContainer, Opts types.ContainerRemoveOptions) error {
	slog.Info("Received remove request", "name", cont.ContainerConfig.Hostname)
	err := n.cli.ContainerRemove(context.Background(), cont.ContID, Opts)
	if err != nil {
		slog.Error("could not remove container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.ContStatus = "removed"
	slog.Info("Container removed", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (n *NodeService) ListContainers(cont *OrchContainer, Opts types.ContainerListOptions) ([]types.Container, error) {
	slog.Info("Received list container request", "name", cont.ContainerConfig.Hostname)
	var containerList []types.Container
	containerList, err := n.cli.ContainerList(context.Background(), Opts)
	if err != nil {
		slog.Error("could not list containers", "name", cont.ContainerConfig.Hostname)
		return nil, err
	}
	return containerList, nil
}

func (n *NodeService) PauseCont(cont *OrchContainer) error {
	slog.Info("Received pause request", "name", cont.ContainerConfig.Hostname)
	err := n.cli.ContainerPause(context.Background(), cont.ContID)
	if err != nil {
		slog.Error("could not pause container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.ContStatus = "paused"
	slog.Info("Container paused", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (n *NodeService) UnpauseCont(cont *OrchContainer) error {
	slog.Info("Received unpause request", "name", cont.ContainerConfig.Hostname)
	err := n.cli.ContainerUnpause(context.Background(), cont.ContID)
	if err != nil {
		slog.Error("could not unpause container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	cont.ContStatus = "running"
	slog.Info("Container unpaused", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (n *NodeService) CopyToCont(cont *OrchContainer, dest string, src io.Reader, Opts types.CopyToContainerOptions) error {
	slog.Info("Received copyTo request", "name", cont.ContainerConfig.Hostname)
	err := n.cli.CopyToContainer(context.Background(), cont.ContID, dest, src, Opts)
	if err != nil {
		slog.Error("could not copyTo container", "name", cont.ContainerConfig.Hostname)
		return err
	}
	slog.Info("copyTo Container", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return nil
}

func (n *NodeService) CopyFromCont(cont *OrchContainer, src string) (io.ReadCloser, error) {
	slog.Info("Received copyFrom request", "name", cont.ContainerConfig.Hostname)
	res, _, err := n.cli.CopyFromContainer(context.Background(), cont.ContID, src)
	if err != nil {
		slog.Error("could not copyFrom container", "name", cont.ContainerConfig.Hostname)
		return nil, err
	}
	slog.Info("copyFrom Container", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return res, nil
}

func (n *NodeService) TopCont(cont *OrchContainer, args []string) (container.ContainerTopOKBody, error) {
	slog.Info("Received Top request", "name", cont.ContainerConfig.Hostname)
	res, err := n.cli.ContainerTop(context.Background(), cont.ContID, args)
	if err != nil {
		slog.Error("could not top container", "name", cont.ContainerConfig.Hostname)
		return res, err
	}
	slog.Info("container top returned", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return res, nil
}

func (n *NodeService) StatCont(cont *OrchContainer, stream bool) (types.ContainerStats, error) {
	slog.Info("Received stat request", "name", cont.ContainerConfig.Hostname)
	res, err := n.cli.ContainerStats(context.Background(), cont.ContID, stream)
	if err != nil {
		slog.Error("could not stat container", "name", cont.ContainerConfig.Hostname)
		return res, err
	}
	slog.Info("container stat returned", "name", cont.ContainerConfig.Hostname, "ID", cont.ContID)
	return res, nil
}
