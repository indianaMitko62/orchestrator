package master

import (
	"fmt"
	"log/slog"
	"net/rpc"

	"github.com/indianaMitko62/orchestrator/src/orchestrator/node"
	"github.com/indianaMitko62/orchestrator/src/orchestrator/noderpc"
)

type MasterService struct {
	nodes map[string]*NodeManager
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
	m.nodes = make(map[string]*NodeManager)
	for _, ns := range nodes {
		m.nodes[ns.Name] = &NodeManager{
			NodeSettings: *ns,
			client:       nil,
		}
	}
	return m
}

func (m *MasterService) ConnectToNodes() error {
	for _, node := range m.nodes {
		err := node.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MasterService) HelloWorld() error {
	for _, node := range m.nodes {
		err := node.HelloWorld()
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

func (n *NodeManager) HelloWorld() error {
	// var reply string
	// arg := "Mitko"
	// err := n.client.Call("NodeServiceRPC.SimpleHello", &arg, &reply)

	var reply noderpc.CreateContReply
	err := n.client.Call("NodeServiceRPC.CreateCont", &noderpc.CreateContArgs{
		Container: &node.ContainerSettings{
			Name: n.Name,
		},
	}, &reply)

	if err != nil {
		return fmt.Errorf("could not call CreateCont: %w", err)
	}
	slog.Info("node returned greeting", "node_name", n.Name, "greeting", reply)
	return nil
}
