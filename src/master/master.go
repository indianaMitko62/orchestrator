package master

import (
	"fmt"
	"net/http"
	"net/rpc"
)

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
	client, err := rpc.DialHTTP("tcp", n.Address)
	if err != nil {
		return fmt.Errorf("could not connect to node's %s RPC service at %s: %w", n.Name, n.Address, err)
	}
	n.Client = client
	return nil
}

func (m *MasterService) Master() {
	http.HandleFunc("/clusterState", m.CS.HandleClusterState)
	http.ListenAndServe(":1986", nil)

	// More Master logic
}
