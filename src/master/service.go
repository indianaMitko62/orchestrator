package master

import (
	"fmt"
	"net/rpc"

	"github.com/docker/docker/client"
)

type MasterService struct {
	cli *client.Client

	CS *ClusterState
}

func NewMasterService(cs *ClusterState) *MasterService {
	m := &MasterService{}
	m.CS = cs
	m.CS.Nodes = make(map[string]*NodeManager)
	for _, nm := range m.CS.Nodes {
		m.CS.Nodes[*nm.Name] = &NodeManager{
			NodeSettings: nm.NodeSettings,
			client:       nil,
		}
	}
	return m
}

func (m *MasterService) ConnectToNodes() error {
	for _, node := range m.CS.Nodes {
		err := node.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *NodeManager) Connect() error {
	client, err := rpc.DialHTTP("tcp", *n.Address)
	if err != nil {
		return fmt.Errorf("could not connect to node's %s RPC service at %s: %w", *n.Name, *n.Address, err)
	}
	n.client = client
	return nil
}
