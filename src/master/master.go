package master

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indianaMitko62/orchestrator/src/cluster"
	"gopkg.in/yaml.v3"
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

	CSToSend, _ := msvc.CS.ToYaml()
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

	var nodeName string // node is alive???
	for name, node := range msvc.CS.Nodes {
		if node.Address == nodeIP {
			nodeName = name
			break
		}
	}

	var outcome cluster.ClusterChangeOutcome
	yamlData, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Error reading YAML data:", err)
	}
	fmt.Println(string(yamlData)) // for testing

	err = yaml.Unmarshal(yamlData, &outcome)
	if err != nil {
		slog.Error("could not unmarshal cluster state yaml", "error", err)

	}
	msvc.ClusterChangeOutcome = &outcome
	slog.Info("Received cluster state request", "node", nodeName, "IP", nodeIP)
	for name, log := range outcome.Logs { // for result
		fmt.Println(name, log)
	}
}

func (m *MasterService) Master() {
	r := mux.NewRouter()
	r.HandleFunc("/clusterState", m.getClusterStateHandler).Methods("GET")
	r.HandleFunc("/clusterState", m.postClusterStateHandler).Methods("POST")
	http.ListenAndServe(":1986", r)

	// More Master logic
}
