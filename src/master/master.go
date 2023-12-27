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

func (m *MasterService) CreateCont() error {
	for _, node := range m.nodes {
		err := node.CreateCont()
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

func (n *NodeManager) CreateCont() error {
	// var reply string
	// arg := "Mitko"
	// err := n.client.Call("NodeServiceRPC.SimpleHello", &arg, &reply)

	contID := ""

	cont := node.Container{
		ContainerConfig: &container.Config{
			Image: "alpine:latest",
			Cmd:   []string{"ping", "localhost"},
		},
		Image_name: "alpine", HostConfig: nil, NetworkingConfig: nil, ContainerName: "cont1", ContID: &contID}

	settings := node.ContainerSettings{Cont: &cont}
	args := noderpc.CreateContArgs{Settings: &settings, Delay: 5}
	var reply noderpc.CreateContReply

	err := n.client.Call("NodeServiceRPC.CreateCont", &args, &reply)

	if err != nil {
		return fmt.Errorf("could not call CreateCont: %w", err)
	}
	slog.Info("Container created", "node_name", n.Name, "ID", reply.ReplyID)
	cont.ContID = &reply.ReplyID
	//gob.Register(types.ContainerStartOptions{})
	settings1 := node.ContainerSettings{Cont: &cont}
	args1 := noderpc.StartContArgs{Settings: &settings1, Delay: 1, Opts: types.ContainerStartOptions{}}
	var reply1 noderpc.CreateContReply
	err = n.client.Call("NodeServiceRPC.StartCont", &args1, &reply1)
	if err != nil {
		return fmt.Errorf("could not call StartCont: %w", err)
	}
	slog.Info("Container started", "node_name", n.Name, "ID", reply1.ReplyID)

	settings = node.ContainerSettings{Cont: &cont}
	args2 := noderpc.StopContArgs{Settings: &settings, Delay: 1, Opts: container.StopOptions{}}
	var reply2 noderpc.CreateContReply
	err = n.client.Call("NodeServiceRPC.StopCont", &args2, &reply2)

	if err != nil {
		return fmt.Errorf("could not call StopCont: %w", err)
	}
	slog.Info("Container stopped", "node_name", n.Name, "ID", reply1.ReplyID)
	return nil
}
