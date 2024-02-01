package master

import (
	"net/http"
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

func (m *MasterService) Master() {
	http.HandleFunc("/clusterState", m.CS.HandleClusterState)
	http.ListenAndServe(":1986", nil)

	// More Master logic
}
