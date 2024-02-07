package master

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

func (m *MasterService) Master() {
	http.HandleFunc("/clusterState", m.HandleClusterState) // to be intergraded with gorrilla muxer
	http.ListenAndServe(":1986", nil)

	// More Master logic
}

func (m *MasterService) HandleClusterState(w http.ResponseWriter, r *http.Request) {
	nodeIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		slog.Error("Error parsing node IP:", err)
		return
	}

	var nodeName string // node is alive???
	for name, node := range m.CS.Nodes {
		if node.Address == nodeIP {
			nodeName = name
			break
		}
	}

	slog.Info("Received cluster state request", "node", nodeName, "IP", nodeIP)

	CSToSend, _ := m.CS.ToYaml()
	fmt.Println("YAML Output:")
	fmt.Println(string(CSToSend))

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(CSToSend)
}
