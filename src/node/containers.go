package node

import (
	"context"
	"io"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type Container struct {
	ContainerConfig  *container.Config
	Image_name       string
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
	ContainerName    string
	ContID           *string
	ContStatus       *string
}

// type ContainerSettings struct { // not needed at the moment, but who knows
// 	Cont *Container
// }

func (n *NodeService) CreateCont(cont *Container) (string, error) {
	slog.Info("Received create request", "name", cont.ContainerName)
	reply, err := n.cli.ContainerCreate(context.Background(),
		cont.ContainerConfig,
		cont.HostConfig,
		cont.NetworkingConfig,
		nil,
		cont.ContainerName)
	if err != nil {
		slog.Error("could not create container", "name", cont.ContainerName)
		return "", err
	}

	cont.ContID = &reply.ID
	*cont.ContStatus = "created"
	slog.Info("Container created", "name", &cont.ContainerName)
	return reply.ID, nil
}

func (n *NodeService) StartCont(Cont *Container, Opts types.ContainerStartOptions) error {
	slog.Info("Received start request", "name", Cont.ContainerName)
	err := n.cli.ContainerStart(context.Background(), *Cont.ContID, Opts)
	if err != nil {
		slog.Error("could not start container", "name", Cont.ContID)
		return err
	}
	*Cont.ContStatus = "running"
	slog.Info("Container started", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) StopCont(Cont *Container, Opts container.StopOptions) error {
	slog.Info("Received stop request", "name", Cont.ContainerName)
	err := n.cli.ContainerStop(context.Background(), *Cont.ContID, Opts)
	if err != nil {
		slog.Error("could not stop container", "name", Cont.ContainerName)
		return err
	}
	*Cont.ContStatus = "stopped"
	slog.Info("Container stopped", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) LogCont(Cont *Container, Opts types.ContainerLogsOptions) (io.ReadCloser, error) {
	slog.Info("Received log request", "name", Cont.ContainerName)
	out, err := n.cli.ContainerLogs(context.Background(), *Cont.ContID, Opts)
	if err != nil {
		slog.Error("could not log container", "name", Cont.ContainerName)
		return nil, err
	}
	slog.Info("Containers Logs returned", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return out, nil
}

func (n *NodeService) RemoveCont(Cont *Container, Opts types.ContainerRemoveOptions) error {
	slog.Info("Received remove request", "name", Cont.ContainerName)
	err := n.cli.ContainerRemove(context.Background(), *Cont.ContID, Opts)
	if err != nil {
		slog.Error("could not remove container", "name", Cont.ContainerName)
		return err
	}
	*Cont.ContStatus = "removed"
	slog.Info("Container removed", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) KillCont(Cont *Container, signal string) error {
	slog.Info("Received kill request", "name", Cont.ContainerName)
	err := n.cli.ContainerKill(context.Background(), *Cont.ContID, signal)
	if err != nil {
		slog.Error("could not kill container", "name", Cont.ContainerName)
		return err
	}
	*Cont.ContStatus = "killed"
	slog.Info("Container killed", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) PauseCont(Cont *Container) error {
	slog.Info("Received pause request", "name", Cont.ContainerName)
	err := n.cli.ContainerPause(context.Background(), *Cont.ContID)
	if err != nil {
		slog.Error("could not pause container", "name", Cont.ContainerName)
		return err
	}
	*Cont.ContStatus = "paused"
	slog.Info("Container paused", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) UnpauseCont(Cont *Container) error {
	slog.Info("Received unpause request", "name", Cont.ContainerName)
	err := n.cli.ContainerUnpause(context.Background(), *Cont.ContID)
	if err != nil {
		slog.Error("could not unpause container", "name", Cont.ContainerName)
		return err
	}
	*Cont.ContStatus = "running"
	slog.Info("Container unpaused", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) CopyToCont(Cont *Container, dest string, src io.Reader, Opts types.CopyToContainerOptions) error {
	slog.Info("Received copyTo request", "name", Cont.ContainerName)
	err := n.cli.CopyToContainer(context.Background(), *Cont.ContID, dest, src, Opts)
	if err != nil {
		slog.Error("could not copyTo container", "name", Cont.ContainerName)
		return err
	}
	slog.Info("copyTo Container", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return nil
}

func (n *NodeService) CopyFromCont(Cont *Container, src string) (io.ReadCloser, error) {
	slog.Info("Received copyFrom request", "name", Cont.ContainerName)
	res, _, err := n.cli.CopyFromContainer(context.Background(), *Cont.ContID, src)
	if err != nil {
		slog.Error("could not copyFrom container", "name", Cont.ContainerName)
		return nil, err
	}
	slog.Info("copyFrom Container", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return res, nil
}

func (n *NodeService) TopCont(Cont *Container, args []string) (container.ContainerTopOKBody, error) {
	slog.Info("Received Top request", "name", Cont.ContainerName)
	res, err := n.cli.ContainerTop(context.Background(), *Cont.ContID, args)
	if err != nil {
		slog.Error("could not top container", "name", Cont.ContainerName)
		return res, err
	}
	slog.Info("container top returned", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return res, nil
}

func (n *NodeService) StatCont(Cont *Container, stream bool) (types.ContainerStats, error) {
	slog.Info("Received stat request", "name", Cont.ContainerName)
	res, err := n.cli.ContainerStats(context.Background(), *Cont.ContID, stream)
	if err != nil {
		slog.Error("could not stat container", "name", Cont.ContainerName)
		return res, err
	}
	slog.Info("container stat returned", "name", Cont.ContainerName, "ID", *Cont.ContID)
	return res, nil
}
