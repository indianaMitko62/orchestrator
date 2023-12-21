package node

import (
	"fmt"
	"log/slog"
)

type NodeService struct {
}

type ContainerSettings struct {
	Name string
}

func NewNodeService() (*NodeService, error) {
	return &NodeService{}, nil
}

func (n *NodeService) CreateCont(settings *ContainerSettings) (string, error) {
	slog.Info("received create request", "name", settings.Name)
	return fmt.Sprintf("Hello, %s", settings.Name), nil
}
