package master

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

func (msvc *MasterService) getClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	nodeIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		slog.Error("Error parsing node IP:", err)
		return
	}

	var nodeName string // node is alive???
	for name, node := range msvc.CS.Nodes {
		if node.Address == nodeIP {
			nodeName = name
			break
		}
	}

	slog.Info("Received cluster state request", "node", nodeName, "IP", nodeIP)

	CSToSend, _ := cluster.ToYaml(msvc.CS)
	fmt.Println("YAML Output:")
	fmt.Println(string(CSToSend))

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(CSToSend)
}

func (msvc *MasterService) postClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	nodeIP, _, err := net.SplitHostPort(r.RemoteAddr) // to be moved to log function. repetition from getClusterStateHandler
	if err != nil {
		slog.Error("Error parsing node IP:", err)
		return
	}
	slog.Info("recieved from", "IP", nodeIP)
	// var nodeName string // node is alive???
	// for name, node := range msvc.CS.Nodes {
	// 	if node.Address == nodeIP {
	// 		nodeName = name
	// 		break
	// 	}
	// }

	logData, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error reading YAML data:", err)
	}
	fmt.Println(string(logData)) // for testing
}

func (msvc *MasterService) postNodeLogsHandler(w http.ResponseWriter, r *http.Request) {
	nodeIP, _, err := net.SplitHostPort(r.RemoteAddr) // to be moved to log function. repetition from getClusterStateHandler
	if err != nil {
		slog.Error("Error parsing node IP:", err)
		return
	}
	slog.Info("recieved from", "IP", nodeIP)
	// var nodeName string // node is alive???
	// for name, node := range msvc.CS.Nodes {
	// 	if node.Address == nodeIP {
	// 		nodeName = name
	// 		break
	// 	}
	// }

	yamlData, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error reading YAML data:", err)
	}
	fmt.Println(string(yamlData)) // for testing
}

func (m *MasterService) Master() {
	r := mux.NewRouter()
	r.HandleFunc("/clusterState", m.getClusterStateHandler).Methods("GET")
	r.HandleFunc("/clusterState", m.postClusterStateHandler).Methods("POST") // move to /logs. POST to /logs from nodes. GET from /logs from CLI
	r.HandleFunc("/nodeStatus", m.postNodeLogsHandler).Methods("POST")
	http.ListenAndServe(":1986", r)

	// More Master logic
}
