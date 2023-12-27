package node

import (
	"context"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type Container struct {
	ContainerConfig  *container.Config
	Image_name       string
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
	ContainerName    string
	ContID           *string
}

type ContainerSettings struct {
	Cont *Container
}

type NodeService struct {
	cli *client.Client
}

func NewNodeService() (*NodeService, error) {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	return &NodeService{cli: cli}, nil
}

func (n *NodeService) CreateCont(settings *ContainerSettings) (string, error) {
	slog.Info("received create request", "name", settings.Cont.ContainerName)

	reply, err := n.cli.ContainerCreate(context.Background(),
		settings.Cont.ContainerConfig,
		settings.Cont.HostConfig,
		settings.Cont.NetworkingConfig,
		nil,
		settings.Cont.ContainerName)

	if err != nil {
		slog.Error("could not create container", "name", settings.Cont.ContainerName)
		return "", err
	}
	settings.Cont.ContID = &reply.ID
	return reply.ID, nil
}

func (n *NodeService) StartCont(settings *ContainerSettings, Opts types.ContainerStartOptions) error {
	err := n.cli.ContainerStart(context.Background(), *settings.Cont.ContID, Opts)
	if err != nil {
		slog.Error("could not start container", "name", settings.Cont.ContID)
		return err
	}
	// out, _ := n.cli.ContainerLogs(context.Background(), *settings.Cont.ContID, types.ContainerLogsOptions{})
	// io.Copy(os.Stdout, out)
	return nil
}

func (n *NodeService) StopCont(settings *ContainerSettings, Opts container.StopOptions) error {
	err := n.cli.ContainerStop(context.Background(), *settings.Cont.ContID, Opts)
	if err != nil {
		slog.Error("could not stop container", "name", settings.Cont.ContainerName)
		return err
	}
	return nil
}
