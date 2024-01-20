package master

import (
	"fmt"
	"log/slog"
	"net/rpc"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"

	"github.com/indianaMitko62/orchestrator/src/node"
	"github.com/indianaMitko62/orchestrator/src/noderpc"
)

type MasterService struct {
	Nodes map[string]*NodeManager
}

type NodeSettings struct {
	Name    string
	Address string
}

type NodeManager struct {
	NodeSettings
	client *rpc.Client
}

func NewMasterService(nodes []*NodeSettings) *MasterService {
	m := &MasterService{}
	m.Nodes = make(map[string]*NodeManager)
	for _, ns := range nodes {
		m.Nodes[ns.Name] = &NodeManager{
			NodeSettings: *ns,
			client:       nil,
		}
	}

	return m
}

func (m *MasterService) ConnectToNodes() error {
	for _, node := range m.Nodes {
		err := node.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MasterService) CreateContOn(node *NodeManager, cont *node.OrchContainer) error {
	err := node.CreateCont(cont)
	if err != nil {
		return err
	}
	return nil
}

func (m *MasterService) StartContOn(node *NodeManager, cont *node.OrchContainer) error {
	err := node.StartCont(cont)
	if err != nil {
		return err
	}
	return nil
}

func (m *MasterService) StopContOn(node *NodeManager, cont *node.OrchContainer) error {
	err := node.StopCont(cont)
	if err != nil {
		return err
	}
	return nil
}

func (m *MasterService) CreateContOnAll(cont *node.OrchContainer) error {
	for _, node := range m.Nodes {
		err := node.CreateCont(cont)
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *NodeManager) Connect() error {
	client, err := rpc.DialHTTP("tcp", n.Address)
	if err != nil {
		return fmt.Errorf("could not connect to node's %s RPC service at %s: %w", n.Name, n.Address, err)
	}
	n.client = client
	return nil
}

func (n *NodeManager) CreateCont(cont *node.OrchContainer) error {
	// var reply string
	// arg := "Mitko"
	// err := n.client.Call("NodeServiceRPC.SimpleHello", &arg, &reply)

	args := noderpc.CreateContArgs{Cont: cont, Delay: 5}
	var reply noderpc.CreateContReply

	err := n.client.Call("NodeServiceRPC.CreateCont", &args, &reply)
	if err != nil {
		return fmt.Errorf("could not call CreateCont: %w", err)
	}

	slog.Info("Container created", "node_name", n.Name, "ID", reply.ReplyID)
	cont.ContID = &reply.ReplyID
	return nil
}

func (n *NodeManager) StartCont(cont *node.OrchContainer) error {

	args := noderpc.StartContArgs{Cont: cont, Delay: 1, Opts: types.ContainerStartOptions{}}
	var reply noderpc.CreateContReply

	err := n.client.Call("NodeServiceRPC.StartCont", &args, &reply)
	if err != nil {
		return fmt.Errorf("could not call StartCont: %w", err)
	}

	slog.Info("Container started", "node_name", n.Name, "ID", *cont.ContID)
	return nil
}

func (n *NodeManager) StopCont(cont *node.OrchContainer) error {

	args := noderpc.StopContArgs{Cont: cont, Delay: 1, Opts: container.StopOptions{}}
	var reply noderpc.CreateContReply

	err := n.client.Call("NodeServiceRPC.StopCont", &args, &reply)
	if err != nil {
		return fmt.Errorf("could not call StopCont: %w", err)
	}

	slog.Info("Container stopped", "node_name", n.Name, "ID", *cont.ContID)
	return nil
}
