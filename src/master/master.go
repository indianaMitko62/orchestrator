package master

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

func (msvc *MasterService) initHTTPServer() {
	r := mux.NewRouter() // separate HTTP server init
	r.HandleFunc("/clusterState", msvc.getClusterStateHandler).Methods("GET")
	r.HandleFunc("/logs", msvc.postLogsHandler).Methods("POST")
	r.HandleFunc("/nodeStatus", msvc.postNodeStatusHandler).Methods("POST")
	go http.ListenAndServe(":1986", r)
}

func (msvc *MasterService) evaluateNodes(inactiveNodeName string) (string, error) {
	bestScore := 1000.0
	var bestNodeName string
	for name, node := range msvc.CS.Nodes {
		if node.Name != inactiveNodeName {
			score := (msvc.NodesStatus[node.Name].CPU + msvc.NodesStatus[node.Name].Memory + msvc.NodesStatus[node.Name].Disc) / 3
			score += float64(len(msvc.NodesStatus[node.Name].CurrentNodeState.Containers))
			if score < bestScore {
				bestScore = score
				bestNodeName = name
			}
		}
	}
	if bestScore == 1000.0 {
		return "", errors.New("no active nodes remaining")
	}
	return bestNodeName, nil
}

func (msvc *MasterService) lostANode(inactiveNodeName string) {
	bestActiveNode, err := msvc.evaluateNodes(inactiveNodeName)
	if err != nil {
		msvc.masterLog.Logger.Error("Could not choose active node", "error", err.Error())
		return
	}
	msvc.masterLog.Logger.Info("Moving containers", "from", inactiveNodeName, "to", bestActiveNode)
	for name, img := range msvc.CS.Nodes[inactiveNodeName].Images {
		msvc.CS.Nodes[bestActiveNode].Images[name] = img
	}
	for name, netw := range msvc.CS.Nodes[inactiveNodeName].Networks {
		msvc.CS.Nodes[bestActiveNode].Networks[name] = netw
	}
	for name, vol := range msvc.CS.Nodes[inactiveNodeName].Volumes {
		msvc.CS.Nodes[bestActiveNode].Volumes[name] = vol
	}
	for name, cont := range msvc.CS.Nodes[inactiveNodeName].Containers {
		msvc.CS.Nodes[bestActiveNode].Containers[name] = cont
	}
}

func (msvc *MasterService) Master() {
	msvc.initHTTPServer()
	for {
		for name, status := range msvc.NodesStatus {
			if time.Since(status.Timestamp) > time.Duration(15*time.Second) {
				msvc.masterLog.Logger.Error("Node inactive", "name", name, "time", time.Since(status.Timestamp))
				status.Active = false
				msvc.lostANode(name)
				msvc.CS.Nodes[name] = cluster.NodeManager{}
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
