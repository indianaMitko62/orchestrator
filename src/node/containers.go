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
	ContainerConfig  *container.Config
	Image_name       string
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
	Name             string
	ContID           *string
	ContStatus       *string
}

/*
NOTES: Most of these functionalities do not need to be accessable remotely.
*/

// type ContainerSettings struct { // not needed at the moment, but who knows
// 	Cont *Container
// }

func (n *NodeService) CreateCont(cont *OrchContainer) (string, error) {
	slog.Info("Received create request", "name", cont.Name)
	reply, err := n.cli.ContainerCreate(context.Background(),
		cont.ContainerConfig,
		cont.HostConfig,
		cont.NetworkingConfig,
		nil,
		cont.Name)
	if err != nil {
		slog.Error("could not create container", "name", cont.Name)
		return "", err
	}

	cont.ContID = &reply.ID
	*cont.ContStatus = "created"
	slog.Info("Container created", "name", &cont.Name)
	return reply.ID, nil
}

func (n *NodeService) StartCont(Cont *OrchContainer, Opts types.ContainerStartOptions) error {
	slog.Info("Received start request", "name", Cont.Name)
	err := n.cli.ContainerStart(context.Background(), *Cont.ContID, Opts)
	if err != nil {
		slog.Error("could not start container", "name", Cont.ContID)
		return err
	}
	*Cont.ContStatus = "running"
	slog.Info("Container started", "name", Cont.Name, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) StopCont(Cont *OrchContainer, Opts container.StopOptions) error {
	slog.Info("Received stop request", "name", Cont.Name)
	err := n.cli.ContainerStop(context.Background(), *Cont.ContID, Opts)
	if err != nil {
		slog.Error("could not stop container", "name", Cont.Name)
		return err
	}
	*Cont.ContStatus = "stopped"
	slog.Info("Container stopped", "name", Cont.Name, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) LogCont(Cont *OrchContainer, Opts types.ContainerLogsOptions) (io.ReadCloser, error) {
	slog.Info("Received log request", "name", Cont.Name)
	out, err := n.cli.ContainerLogs(context.Background(), *Cont.ContID, Opts)
	if err != nil {
		slog.Error("could not log container", "name", Cont.Name)
		return nil, err
	}
	slog.Info("Containers Logs returned", "name", Cont.Name, "ID", *Cont.ContID)
	return out, nil
}

func (n *NodeService) KillCont(Cont *OrchContainer, signal string) error {
	slog.Info("Received kill request", "name", Cont.Name)
	err := n.cli.ContainerKill(context.Background(), *Cont.ContID, signal)
	if err != nil {
		slog.Error("could not kill container", "name", Cont.Name)
		return err
	}
	*Cont.ContStatus = "killed"
	slog.Info("Container killed", "name", Cont.Name, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) RemoveCont(Cont *OrchContainer, Opts types.ContainerRemoveOptions) error {
	slog.Info("Received remove request", "name", Cont.Name)
	err := n.cli.ContainerRemove(context.Background(), *Cont.ContID, Opts)
	if err != nil {
		slog.Error("could not remove container", "name", Cont.Name)
		return err
	}
	*Cont.ContStatus = "removed"
	slog.Info("Container removed", "name", Cont.Name, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) ListContainers(Cont *OrchContainer, Opts types.ContainerListOptions) ([]types.Container, error) {
	slog.Info("Received list container request", "name", Cont.Name)
	var containerList []types.Container
	containerList, err := n.cli.ContainerList(context.Background(), Opts)
	if err != nil {
		slog.Error("could not list containers", "name", Cont.Name)
		return nil, err
	}
	return containerList, nil
}

func (n *NodeService) PauseCont(Cont *OrchContainer) error {
	slog.Info("Received pause request", "name", Cont.Name)
	err := n.cli.ContainerPause(context.Background(), *Cont.ContID)
	if err != nil {
		slog.Error("could not pause container", "name", Cont.Name)
		return err
	}
	*Cont.ContStatus = "paused"
	slog.Info("Container paused", "name", Cont.Name, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) UnpauseCont(Cont *OrchContainer) error {
	slog.Info("Received unpause request", "name", Cont.Name)
	err := n.cli.ContainerUnpause(context.Background(), *Cont.ContID)
	if err != nil {
		slog.Error("could not unpause container", "name", Cont.Name)
		return err
	}
	*Cont.ContStatus = "running"
	slog.Info("Container unpaused", "name", Cont.Name, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) CopyToCont(Cont *OrchContainer, dest string, src io.Reader, Opts types.CopyToContainerOptions) error {
	slog.Info("Received copyTo request", "name", Cont.Name)
	err := n.cli.CopyToContainer(context.Background(), *Cont.ContID, dest, src, Opts)
	if err != nil {
		slog.Error("could not copyTo container", "name", Cont.Name)
		return err
	}
	slog.Info("copyTo Container", "name", Cont.Name, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) CopyFromCont(Cont *OrchContainer, src string) (io.ReadCloser, error) {
	slog.Info("Received copyFrom request", "name", Cont.Name)
	res, _, err := n.cli.CopyFromContainer(context.Background(), *Cont.ContID, src)
	if err != nil {
		slog.Error("could not copyFrom container", "name", Cont.Name)
		return nil, err
	}
	slog.Info("copyFrom Container", "name", Cont.Name, "ID", *Cont.ContID)
	return res, nil
}

func (n *NodeService) TopCont(Cont *OrchContainer, args []string) (container.ContainerTopOKBody, error) {
	slog.Info("Received Top request", "name", Cont.Name)
	res, err := n.cli.ContainerTop(context.Background(), *Cont.ContID, args)
	if err != nil {
		slog.Error("could not top container", "name", Cont.Name)
		return res, err
	}
	slog.Info("container top returned", "name", Cont.Name, "ID", *Cont.ContID)
	return res, nil
}

func (n *NodeService) StatCont(Cont *OrchContainer, stream bool) (types.ContainerStats, error) {
	slog.Info("Received stat request", "name", Cont.Name)
	res, err := n.cli.ContainerStats(context.Background(), *Cont.ContID, stream)
	if err != nil {
		slog.Error("could not stat container", "name", Cont.Name)
		return res, err
	}
	slog.Info("container stat returned", "name", Cont.Name, "ID", *Cont.ContID)
	return res, nil
}
