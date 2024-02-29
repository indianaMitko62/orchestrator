package master

import (
	"fmt"
	"io"
	"net/http"

	"github.com/indianaMitko62/orchestrator/src/cluster"
	"gopkg.in/yaml.v3"
)

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
	if msvc.NodesStatusLogs[nodeName] == nil { // separate node Status logging functions
		log := cluster.NewLog(msvc.LogsPath + nodeName + "StatusLogs")
		if log != nil {
			msvc.NodesStatusLogs[nodeName] = log
		}
	}
	var unhealthyCnt = 0
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
			msvc.NodesStatusLogs[nodeName].Logger.Info("Node Status", "Node", nodeName, "CPU", nodeStatus.CPU, "Memory", nodeStatus.Memory, "Disk", nodeStatus.Disk,
				"HealthyContainers", healthyContLog, "UnhealthyContainers", unhealthyContLog)
		} else {
			msvc.NodesStatusLogs[nodeName].Logger.Info("Node Status", "Node", nodeName, "CPU", nodeStatus.CPU, "Memory", nodeStatus.Memory, "Disk", nodeStatus.Disk,
				"HealthyContainers", healthyContLog)
		}
	}
}
