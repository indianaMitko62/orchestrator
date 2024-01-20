package noderpc

import (
	"fmt"

	"github.com/indianaMitko62/orchestrator/src/node"
)

type NodeServiceRPC struct {
	service *node.NodeService
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
