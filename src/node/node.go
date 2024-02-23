package node

import (
	"fmt"
	"time"

	"github.com/indianaMitko62/orchestrator/src/cluster"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
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
	err := nsvc.sendLogs(nsvc.MasterAddress+nsvc.LogsPath, nsvc.clusterChangeLog)
	return err
}

func (nsvc *NodeService) applyChanges() error {
	nsvc.nodeLog.Logger.Info("finding differences")
	if nsvc.changeContainers() || nsvc.changeVolumes() || nsvc.changeNetworks() {
		return nsvc.sendLogs(nsvc.MasterAddress+nsvc.LogsPath, nsvc.clusterChangeLog)
	} else {
		nsvc.nodeLog.Logger.Info("No changes in cluster")
	}
	return nil
}

func (nsvc *NodeService) inspectCluster() {
	if nsvc.CurrentNodeState == nil {
		nsvc.nodeLog.Logger.Error("No Current Node State")
		return
	}

	for _, cont := range nsvc.CurrentNodeState.Containers {
		if cont.CurrentStatus == "running" {
			nsvc.inspectContainer(cont)
		}
	}
	percent, _ := cpu.Percent(time.Second, false)
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	ns := cluster.NodeStatus{
		CPU:              percent[0], // add these
		Memory:           memInfo.UsedPercent,
		Disc:             40,
		CurrentNodeState: *nsvc.CurrentNodeState,
		Active:           true,
		Timestamp:        time.Now(),
	}
	nsvc.SendNodeStatus(nsvc.MasterAddress+nsvc.NodeStatusPath, &ns)
}

func (nsvc *NodeService) Node() error {
	nsvc.MasterAddress = "http://" + nsvc.MasterAddress + nsvc.MasterPort
	for {
		err := nsvc.getClusterState(nsvc.MasterAddress + nsvc.ClusterStatePath)
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
		// time.Sleep(time.Duration(5-time.Now().Second()%5) * time.Second)
		time.Sleep(40 * time.Second) // for node inactivity simulation
		fmt.Print("\n\n\n")
	}
}
