package master

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/indianaMitko62/orchestrator/src/cluster"
	"gopkg.in/yaml.v3"
)

func (msvc *MasterService) getClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	nodeName := r.Header.Get("nodeName")
	msvc.masterLog.Logger.Info("Received GET on /clusterState from", "node", nodeName, "IP", r.RemoteAddr)

	CSToSend, _ := cluster.ToYaml(msvc.CS)
	fmt.Println("YAML Output:") // for testing
	fmt.Println(string(CSToSend))

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(CSToSend)
}

func (msvc *MasterService) postLogsHandler(w http.ResponseWriter, r *http.Request) {
	nodeName := r.Header.Get("nodeName")
	msvc.masterLog.Logger.Info("recieved POST on /logs from", "name", nodeName, "IP", r.RemoteAddr)

	logData, err := io.ReadAll(r.Body)
	if err != nil {
		msvc.masterLog.Logger.Error("Error reading YAML data:", err)
	}
	fmt.Println(string(logData)) // for testing
	f, err := os.OpenFile(msvc.LogsPath+nodeName+"Logs", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		msvc.masterLog.Logger.Warn("Could not open file. Trying to create it", "name", msvc.LogsPath+nodeName+"Logs")
		_, err := os.Create(msvc.LogsPath + nodeName + "Logs")
		if err != nil {
			msvc.masterLog.Logger.Error("Could not create file", "name", msvc.LogsPath+nodeName+"Logs", "error", err)
			return
		}
		f, _ = os.OpenFile(msvc.LogsPath+nodeName+"Logs", os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			msvc.masterLog.Logger.Error("Could not open file", "name", msvc.LogsPath+nodeName+"Logs", "error", err)
			return
		}
	}
	f.Write(logData)
	f.Close()
	msvc.masterLog.Logger.Info("Logs written", "file", msvc.LogsPath+nodeName+"Logs")
}

func (msvc *MasterService) postNodeStatusHandler(w http.ResponseWriter, r *http.Request) {
	nodeName := r.Header.Get("nodeName")
	msvc.masterLog.Logger.Info("recieved POST on /nodeStatus from", "name", nodeName, "IP", r.RemoteAddr)
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
	var unhealthyCnt = 0 // separate node Status logging functions
	var healthyCnt = 0
	var startingCnt = 0
	unhealthyConts := ""
	var unhealthyContLog string
	var healthyContLog string
	for name, cont := range nodeStatus.CurrentNodeState.Containers {
		if cont.CurrHealth == "unhealthy" {
			unhealthyCnt++
			unhealthyConts += name + ", "
		}
		if cont.CurrHealth == "healthy" {
			healthyCnt++
		}
		if cont.CurrHealth == "starting" {
			startingCnt++
		}
	}
	if healthyCnt+unhealthyCnt == 0 {
		if startingCnt == 0 {
			msvc.NodesStatusLogs[nodeName].Logger.Info("No containers on node")
		} else {
			msvc.NodesStatusLogs[nodeName].Logger.Info("All Containers starting")
		}
	} else {
		healthyPercent := healthyCnt / (healthyCnt + unhealthyCnt) * 100
		healthyContLog = "Healthy containers: " + fmt.Sprintf("%f", float32(healthyPercent))
		fmt.Println(healthyCnt + unhealthyCnt)
		if unhealthyConts != "" {
			unhealthyContLog = "Unhealthy containers: " + unhealthyConts
			msvc.NodesStatusLogs[nodeName].Logger.Info("Node Status", "Node", nodeName, "CPU", nodeStatus.CPU, "Memory", nodeStatus.Memory, "Disk", nodeStatus.Disc,
				"HealthyContainers", healthyContLog, "UnhealthyContainers", unhealthyContLog)
		} else {
			msvc.NodesStatusLogs[nodeName].Logger.Info("Node Status", "Node", nodeName, "CPU", nodeStatus.CPU, "Memory", nodeStatus.Memory, "Disk", nodeStatus.Disc,
				"HealthyContainers", healthyContLog)
		}
	}
}

func (msvc *MasterService) Master() {
	r := mux.NewRouter() // separate HTTP server init
	r.HandleFunc("/clusterState", msvc.getClusterStateHandler).Methods("GET")
	r.HandleFunc("/logs", msvc.postLogsHandler).Methods("POST")
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
