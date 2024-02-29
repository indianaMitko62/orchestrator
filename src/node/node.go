package node

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/indianaMitko62/orchestrator/src/cluster"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

/*
TODO:	CLI basics... ????
*/

func (nsvc *NodeService) stopAllContainersOnMachine() error {
	nsvc.nodeLog.Logger.Info("Stopping all containers on node")
	containers, err := nsvc.cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}
	for _, cont := range containers {
		if err := nsvc.cli.ContainerStop(context.Background(), cont.ID, container.StopOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (nsvc *NodeService) initCluster() error {
	nsvc.CurrentNodeState = cluster.NewNodeState()
	fmt.Println()
	nsvc.stopAllContainersOnMachine()
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
	err := nsvc.sendLogs(nsvc.MasterAddress+nsvc.LogsEndpoint, nsvc.clusterChangeLog)
	return err
}

func (nsvc *NodeService) applyChanges() error {
	nsvc.nodeLog.Logger.Info("finding differences")
	if nsvc.changeContainers() || nsvc.changeVolumes() || nsvc.changeNetworks() {
		return nsvc.sendLogs(nsvc.MasterAddress+nsvc.LogsEndpoint, nsvc.clusterChangeLog)
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
			nsvc.getContHealth(cont)
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
		Disk:             40,
		CurrentNodeState: *nsvc.CurrentNodeState,
		Active:           true,
		Timestamp:        time.Now(),
	}
	nsvc.SendNodeStatus(nsvc.MasterAddress+nsvc.NodeStatusEndpoint, &ns)
}

func (nsvc *NodeService) Node() error {
	nsvc.MasterAddress = "http://" + nsvc.MasterAddress + nsvc.MasterPort
	for {
		err := nsvc.getClusterState(nsvc.MasterAddress + nsvc.ClusterStateEndpoint)
		if err != nil {
			nsvc.nodeLog.Logger.Error("could not get cluster data", "error", err)
		} else {
			if nsvc.DesiredNodeState != nil {
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
			} else {
				nsvc.nodeLog.Logger.Info("No desired node state")
			}
		}
		nsvc.inspectCluster()
		nsvc.nodeLog.Logger.Info("Main Node process sleeping...") // not to be logged everytime. Stays for now for development purposes
		time.Sleep(time.Duration(5-time.Now().Second()%5) * time.Second)
		//time.Sleep(40 * time.Second) // for node inactivity simulation
		fmt.Print("\n\n\n")
	}
}
