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
	r.HandleFunc("/clusterState", msvc.postClusterStateHandler).Methods("POST")
	r.HandleFunc("/logs", msvc.postLogsHandler).Methods("POST")
	r.HandleFunc("/nodeStatus", msvc.postNodeStatusHandler).Methods("POST")
	go http.ListenAndServe(msvc.HTTPServerPort, r)
}

func (msvc *MasterService) evaluateNodes(inactiveNodeName string) (string, error) {
	bestScore := 1000.0
	var bestNodeName string
	for name := range msvc.CS.Nodes {
		fmt.Println(name, msvc.NodesStatus[name].Operating)
		if name != inactiveNodeName && msvc.NodesStatus[name].Operating {
			score := (msvc.NodesStatus[name].CPU + msvc.NodesStatus[name].Memory + msvc.NodesStatus[name].Disk) / 3
			score += float64(len(msvc.NodesStatus[name].CurrentNodeState.Containers))
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
		fmt.Println(name)
	}
	for name, netw := range msvc.CS.Nodes[inactiveNodeName].Networks {
		msvc.CS.Nodes[bestActiveNode].Networks[name] = netw
		fmt.Println(name)
	}
	for name, vol := range msvc.CS.Nodes[inactiveNodeName].Volumes {
		msvc.CS.Nodes[bestActiveNode].Volumes[name] = vol
		fmt.Println(name)
	}
	for name, cont := range msvc.CS.Nodes[inactiveNodeName].Containers {
		msvc.CS.Nodes[bestActiveNode].Containers[name] = cont
		fmt.Println(name)
	}
	node := msvc.CS.Nodes[bestActiveNode]
	node.Active = true
	msvc.CS.Nodes[bestActiveNode] = node
}

func (msvc *MasterService) Master() {
	msvc.initHTTPServer()
	for {
		for name, status := range msvc.NodesStatus {
			if time.Since(status.Timestamp) > time.Duration(15*time.Second) && len(status.CurrentNodeState.Containers) > 0 && status.Operating && msvc.CS.Nodes[name].Active {
				msvc.masterLog.Logger.Error("Node inactive", "name", name, "time", time.Since(status.Timestamp))
				status.Operating = false
				node := msvc.CS.Nodes[name]
				node.Active = false
				msvc.CS.Nodes[name] = node
				msvc.lostANode(name)
				msvc.CS.Nodes[name] = cluster.NodeState{}
			} else {
				msvc.masterLog.Logger.Info("Node", "name", name, "operating", status.Operating, "time", time.Since(status.Timestamp))
			}
		}
		msvc.masterLog.Logger.Info("Main Master process sleeping...") // not to be logged everytime. Stays for now for development purposes
		time.Sleep(time.Duration(5-time.Now().Second()%5) * time.Second)
		fmt.Print("\n\n\n")
	}
}
