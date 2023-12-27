package noderpc

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/indianaMitko62/orchestrator/src/orchestrator/node"
)

type NodeServiceRPC struct {
	service *node.NodeService
}

type CreateContArgs struct {
	// Cont     *node.Container
	Settings *node.ContainerSettings
	Delay    int
}

type StartContArgs struct {
	// Cont     *node.Container
	Settings *node.ContainerSettings
	Delay    int
	Opts     types.ContainerStartOptions
}

type StopContArgs struct {
	// Cont     *node.Container
	Settings *node.ContainerSettings
	Delay    int
	Opts     container.StopOptions
}

type CreateContReply struct {
	ReplyID string
}

func NewNodeServiceRPC() (*NodeServiceRPC, error) {
	svc, err := node.NewNodeService()
	if err != nil {
		return nil, err
	}
	return &NodeServiceRPC{
		service: svc,
	}, nil
}

func (n *NodeServiceRPC) SimpleHello(args *string, reply *string) error {
	*reply = fmt.Sprintf("Hello, %s", *args)
	return nil
}

func (n *NodeServiceRPC) CreateCont(args *CreateContArgs, reply *CreateContReply) error {
	replyID, err := n.service.CreateCont(args.Settings)
	if err != nil {
		return err
	}

	*reply = CreateContReply{
		ReplyID: replyID,
	}
	return nil
}

func (n *NodeServiceRPC) StartCont(args *StartContArgs, reply *CreateContReply) error {
	err := n.service.StartCont(args.Settings, args.Opts)
	if err != nil {
		return err
	}
	return nil
}

func (n *NodeServiceRPC) StopCont(args *StopContArgs, reply *CreateContReply) error {
	err := n.service.StopCont(args.Settings, args.Opts)
	if err != nil {
		return err
	}
	return nil
}
