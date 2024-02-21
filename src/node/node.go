package node

import (
	"fmt"
	"time"

	"github.com/indianaMitko62/orchestrator/src/cluster"
)

/*
TODO: functions managing overall node performance and loading(cpu, memory, disk) and overall node logic
*/

func (nsvc *NodeService) initCluster() error {
	nsvc.CurrentNodeState = cluster.NewNodeState()
	fmt.Println()

	for _, img := range nsvc.DesiredNodeState.Images {
		nsvc.deployNewImage(img)
	}
	fmt.Println()

	for _, netw := range nsvc.DesiredNodeState.Networks {
		nsvc.deployNetwork(netw)
	}
	fmt.Println()

	for _, vol := range nsvc.DesiredNodeState.Volumes {
		nsvc.deployVolume(vol)
	}
	fmt.Println()

	for _, cont := range nsvc.DesiredNodeState.Containers {
		nsvc.deployContainer(cont)
	}
	fmt.Println()
	nsvc.sendLogs(nsvc.MasterAddress+nsvc.ClusterStatePath, nsvc.clusterChangeLog)
	return nil
}

func (nsvc *NodeService) applyChanges() error {
	nsvc.nodeLog.Logger.Info("finding differences")
	if nsvc.changeContainers() || nsvc.changeVolumes() || nsvc.changeNetworks() {
		nsvc.sendLogs(nsvc.MasterAddress+nsvc.ClusterStatePath, nsvc.clusterChangeLog)
	} else {
		nsvc.nodeLog.Logger.Info("No changes in cluster")
	}
	return nil
}

func (nsvc *NodeService) inspectCluster() {
	for _, cont := range nsvc.CurrentNodeState.Containers {
		nsvc.inspectContainer(cont)
	}
	ns := cluster.NodeStatus{
		CPU:              50, // add these
		Memory:           10,
		Disc:             40,
		CurrentNodeState: *nsvc.CurrentNodeState,
	}
	nsvc.SendNodeStatus(nsvc.MasterAddress+nsvc.NodeStatusPath, &ns)
}

func (nsvc *NodeService) Node() error {
	nsvc.MasterAddress = "http://" + nsvc.MasterAddress + nsvc.Port
	clusterStateURL := nsvc.MasterAddress + nsvc.ClusterStatePath // move these logs to /logs - to be separeted in master for different nodes
	for {
		err := nsvc.getClusterState(clusterStateURL)
		if err != nil {
			nsvc.nodeLog.Logger.Error("could not get cluster data", "error", err)
		} else {
			if nsvc.CurrentNodeState == nil {
				nsvc.nodeLog.Logger.Info("No current node state")
				err := nsvc.initCluster()
				if err != nil {
					nsvc.nodeLog.Logger.Error("Could not init cluster")
				}
			} else {
				nsvc.nodeLog.Logger.Info("Present current node state")
				nsvc.applyChanges()
			}
		}
		nsvc.inspectCluster()
		nsvc.nodeLog.Logger.Info("Main Node process sleeping...") // not to be logged everytime. Stays for now for development purposes
		time.Sleep(time.Duration(5-time.Now().Second()%5) * time.Second)
		fmt.Print("\n\n\n")
	}
}
