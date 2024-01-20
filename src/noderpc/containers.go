package noderpc

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/indianaMitko62/orchestrator/src/node"
)

// type ContInfo struct {
// }

type CreateContArgs struct {
	Cont  *node.OrchContainer
	Delay int // example value for context settings later on
}

type StartContArgs struct {
	Cont  *node.OrchContainer
	Delay int // example value for context settings later on
	Opts  types.ContainerStartOptions
}

type StopContArgs struct {
	Cont  *node.OrchContainer
	Delay int // example value for context settings later on
	Opts  container.StopOptions
}

type CreateContReply struct {
	ReplyID string
}

/*
TODO: Most of these functionalities do not need to be accessable remotely. To be refactored...
*/

func (n *NodeServiceRPC) CreateCont(args *CreateContArgs, reply *CreateContReply) error {
	replyID, err := n.service.CreateCont(args.Cont)
	if err != nil {
		return err
	}

	*reply = CreateContReply{
		ReplyID: replyID,
	}
	return nil
}

func (n *NodeServiceRPC) StartCont(args *StartContArgs, reply *CreateContReply) error {
	err := n.service.StartCont(args.Cont, args.Opts)
	if err != nil {
		return err
	}
	return nil
}

func (n *NodeServiceRPC) StopCont(args *StopContArgs, reply *CreateContReply) error {
	err := n.service.StopCont(args.Cont, args.Opts)
	if err != nil {
		return err
	}
	return nil
}
