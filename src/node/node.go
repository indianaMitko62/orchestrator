package node

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/indianaMitko62/orchestrator/src/cluster"
)

/*
TODO: functions managing overall node performance and loading(cpu, memory, disk) and overall node logic
*/

func (nsvc *NodeService) InitCluster() error { // status change. Refactor function needed
	nsvc.CurrentNodeState = cluster.NewNodeState()
	fmt.Println()

	for _, img := range nsvc.DesiredNodeState.Images {
		nsvc.DeployNewImage(img)
	}
	fmt.Println()

	for _, netw := range nsvc.DesiredNodeState.Networks {
		nsvc.DeployNewNetwork(netw)
	}
	fmt.Println()

	for _, vol := range nsvc.DesiredNodeState.Volumes {
		nsvc.DeployNewVolume(vol)
	}
	fmt.Println()

	for _, cont := range nsvc.DesiredNodeState.Containers {
		nsvc.DeployNewContainer(cont)
	}
	fmt.Println()
	nsvc.postClusterChangeOutcome(nsvc.MasterAddress + "/clusterState")
	return nil
}

func (nsvc *NodeService) postClusterChangeOutcome(URL string) {
	req, err := http.NewRequest(http.MethodPost, URL, nsvc.clusterChangeLog.logReader)
	if err != nil {
		nsvc.nodeLog.Logger.Error("Could not create POST request", "URL", URL)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		nsvc.nodeLog.Logger.Error("Could not send POST request")
	}

	if resp.StatusCode == http.StatusOK {
		nsvc.nodeLog.Logger.Info("Cluster Change Outcome logs send successfully")
	}
	file, _ := os.Open("./clusterChangeLog")
	file.Seek(-1, io.SeekEnd)
	nsvc.clusterChangeLog.logReader = file
}

func (nsvc *NodeService) findDifferences() error {
	nsvc.nodeLog.Logger.Info("finding differences")
	if nsvc.changeContainers() || nsvc.changeVolumes() || nsvc.changeNetworks() {
		nsvc.postClusterChangeOutcome(nsvc.MasterAddress + "/clusterState")
	} else {
		nsvc.nodeLog.Logger.Info("No changes in cluster")
	}
	return nil
}

func (nsvc *NodeService) Node() error {
	nsvc.MasterAddress = "http://localhost:1986" //harcoded for now
	clusterStateURL := nsvc.MasterAddress + "/clusterState"
	for {
		recievedClusterState, err := cluster.GetClusterState(clusterStateURL)
		if err != nil {
			nsvc.nodeLog.Logger.Error("could not get cluster data", "error", err)
		} else {
			nsvc.DesiredNodeState = &recievedClusterState.Nodes[nsvc.Name].NodeState
			if nsvc.CurrentNodeState == nil {
				nsvc.nodeLog.Logger.Info("No current node state")
				err := nsvc.InitCluster()
				if err != nil {
					nsvc.nodeLog.Logger.Error("Could not init cluster")
				}
			} else {
				nsvc.nodeLog.Logger.Info("Present current node state")
				nsvc.findDifferences()
			}
		}
		nsvc.nodeLog.Logger.Info("Main Node process sleeping...") // to be moved to different logger, not the one send to /clusterState
		time.Sleep(time.Second * 5)
		fmt.Print("\n\n\n")
	}
}
