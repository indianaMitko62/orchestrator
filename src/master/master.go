package master

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

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
		msvc.masterLog.Logger.Error("Error parsing node IP:", err)
		return
	}
	msvc.masterLog.Logger.Info("recieved from", "IP", nodeIP)

	nodeName := msvc.getNodeNameByIP(nodeIP)
	msvc.masterLog.Logger.Info("Received GET on /clusterState from", "node", nodeName, "IP", nodeIP)

	CSToSend, _ := cluster.ToYaml(msvc.CS)
	fmt.Println("YAML Output:")
	fmt.Println(string(CSToSend))

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(CSToSend)
}

func (msvc *MasterService) postLogsHandler(w http.ResponseWriter, r *http.Request) {
	nodeIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		msvc.masterLog.Logger.Error("Error parsing node IP:", err)
		return
	}

	nodeName := msvc.getNodeNameByIP(nodeIP)
	msvc.masterLog.Logger.Info("recieved POST on /clusterState from", "name", nodeName, "IP", nodeIP)

	logData, err := io.ReadAll(r.Body)
	if err != nil {
		msvc.masterLog.Logger.Error("Error reading YAML data:", err)
	}
	fmt.Println(string(logData)) // for testing
}

func (msvc *MasterService) postNodeStatusHandler(w http.ResponseWriter, r *http.Request) {
	nodeIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		msvc.masterLog.Logger.Error("Error parsing node IP:", err)
		return
	}

	nodeName := msvc.getNodeNameByIP(nodeIP)
	msvc.masterLog.Logger.Info("recieved POST on /nodeStatus from", "name", nodeName, "IP", nodeIP)

	var nodeStatus cluster.NodeStatus
	yamlData, err := io.ReadAll(r.Body)
	if err != nil {
		msvc.masterLog.Logger.Error("Error reading YAML data:", err)
	}
	// fmt.Println(string(yamlData)) // for testing
	yaml.Unmarshal(yamlData, &nodeStatus)
	msvc.NodesStatus[nodeName] = nodeStatus
	if msvc.NodesStatusLogs[nodeName] == nil {
		log := cluster.NewLog("./logs/masterLogs/" + nodeName + "StatusLogs")
		if log != nil {
			msvc.NodesStatusLogs[nodeName] = log
		}
	}
	msvc.NodesStatusLogs[nodeName].Logger.Info("Node Status", "Node", nodeName, "CPU", nodeStatus.CPU, "Memory", nodeStatus.Memory, "Disk", nodeStatus.Disc)
	// To add container stat to logging ^.
	// If there is a unhealthy containers, "Healthy containers: N%: ContainerName unhealthy. Restarted."
	// If none - "Healthy containers percent: 100%"
}

func (msvc *MasterService) Master() {
	r := mux.NewRouter() // separate HTTP server init
	r.HandleFunc("/clusterState", msvc.getClusterStateHandler).Methods("GET")
	r.HandleFunc("/logs", msvc.postLogsHandler).Methods("POST") // move to /logs. POST to /logs from nodes. GET from /logs from CLI
	r.HandleFunc("/nodeStatus", msvc.postNodeStatusHandler).Methods("POST")
	go http.ListenAndServe(":1986", r)

	for {
		for name, status := range msvc.NodesStatus {
			if time.Since(status.Timestamp) > time.Duration(15*time.Second) {
				msvc.masterLog.Logger.Error("Node inactive", "name", name, "time", time.Since(status.Timestamp))
			} else {
				msvc.masterLog.Logger.Info("Node active", "name", name, "time", time.Since(status.Timestamp))
			}
		}
		msvc.masterLog.Logger.Info("Main Master process sleeping...") // not to be logged everytime. Stays for now for development purposes
		time.Sleep(time.Duration(5-time.Now().Second()%5) * time.Second)
		fmt.Print("\n\n\n")

	}
	// More Master logic
}
