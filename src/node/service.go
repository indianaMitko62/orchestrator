package node

import (
	"github.com/docker/docker/client"
)

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
