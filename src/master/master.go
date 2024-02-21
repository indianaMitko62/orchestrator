package master

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indianaMitko62/orchestrator/src/cluster"
	"gopkg.in/yaml.v3"
)

func (msvc *MasterService) getNodeNameByIP(nodeIP string) string {
	for name, node := range msvc.CS.Nodes {
		if node.Address == nodeIP {
			return name
		}
	}
	return ""
}

func (msvc *MasterService) getClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	nodeIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		msvc.MasterLog.Logger.Error("Error parsing node IP:", err)
		return
	}
	msvc.MasterLog.Logger.Info("recieved from", "IP", nodeIP)

	nodeName := msvc.getNodeNameByIP(nodeIP)
	msvc.MasterLog.Logger.Info("Received cluster state request", "node", nodeName, "IP", nodeIP)

	CSToSend, _ := cluster.ToYaml(msvc.CS)
	fmt.Println("YAML Output:")
	fmt.Println(string(CSToSend))

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(CSToSend)
}

func (msvc *MasterService) postClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	nodeIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		msvc.MasterLog.Logger.Error("Error parsing node IP:", err)
		return
	}
	msvc.MasterLog.Logger.Info("recieved from", "IP", nodeIP)

	logData, err := io.ReadAll(r.Body)
	if err != nil {
		msvc.MasterLog.Logger.Error("Error reading YAML data:", err)
	}
	fmt.Println(string(logData)) // for testing
}

func (msvc *MasterService) postNodeStatusHandler(w http.ResponseWriter, r *http.Request) {
	nodeIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		msvc.MasterLog.Logger.Error("Error parsing node IP:", err)
		return
	}
	msvc.MasterLog.Logger.Info("recieved from", "IP", nodeIP)
	nodeName := msvc.getNodeNameByIP(nodeIP)

	var nodeStatus cluster.NodeStatus
	yamlData, err := io.ReadAll(r.Body)
	if err != nil {
		msvc.MasterLog.Logger.Error("Error reading YAML data:", err)
	}
	fmt.Println(string(yamlData)) // for testing
	yaml.Unmarshal(yamlData, &nodeStatus)
	msvc.NodesStatus[nodeName] = nodeStatus
}

func (m *MasterService) Master() {
	r := mux.NewRouter() // separate HTTP server init
	r.HandleFunc("/clusterState", m.getClusterStateHandler).Methods("GET")
	r.HandleFunc("/clusterState", m.postClusterStateHandler).Methods("POST") // move to /logs. POST to /logs from nodes. GET from /logs from CLI
	r.HandleFunc("/nodeStatus", m.postNodeStatusHandler).Methods("POST")
	http.ListenAndServe(":1986", r)

	// More Master logic
}
