package noderpc

import (
	"fmt"

	"github.com/indianaMitko62/orchestrator/src/orchestrator/node"
)

type NodeServiceRPC struct {
	service *node.NodeService
}

type CreateContArgs struct {
	Container *node.ContainerSettings
}

type CreateContReply struct {
	Greeting string
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
	greeting, err := n.service.CreateCont(args.Container)
	if err != nil {
		return err
	}

	*reply = CreateContReply{
		Greeting: greeting,
	}
	return nil
}

type StartContArgs struct {
	Container *node.ContainerSettings
}

func (n *NodeServiceRPC) StartCont(args *StartContArgs, reply *StartContArgs) error {
	return nil
}
