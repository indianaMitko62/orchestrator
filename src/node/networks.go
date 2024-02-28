package node

import (
	"reflect"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

func (nsvc *NodeService) deployNetwork(netw *cluster.OrchNetwork) {
	netw.Cli = nsvc.cli
	_, err := netw.CreateNet(netw.NetworkConfig)
	if err != nil {
		err = nsvc.handleDuplicateNetworks(netw)
	}
	if netw.DesiredStatus == netw.CurrentStatus && err == nil {
		nsvc.CurrentNodeState.Networks[netw.Name] = netw
		nsvc.clusterChangeLog.Logger.Info("Network successfully created", "name", nsvc.CurrentNodeState.Networks[netw.Name].Name)
	} else if err != nil {
		nsvc.clusterChangeLog.Logger.Info("Could not create network", "name", netw.Name, "err", err.Error())
	}
}

func (nsvc *NodeService) handleDuplicateNetworks(newNet *cluster.OrchNetwork) error {
	nsvc.nodeLog.Logger.Info("Trying to remove duplicate network if exists", "name", newNet.Name)
	newNet.RemoveNet()

	nsvc.nodeLog.Logger.Info("Trying to create network again", "name", newNet.Name)
	_, err := newNet.CreateNet(newNet.NetworkConfig)
	if err != nil {
		nsvc.nodeLog.Logger.Error("Second attempt for network creation failed. Aborting...", "name", newNet.Name)
		return err
	}
	return nil
}

func (nsvc *NodeService) stopContainersOnNetwork(netw cluster.OrchNetwork) error {
	netwData, err := netw.InspectNet(types.NetworkInspectOptions{})
	if err != nil {
		return err
	}
	containers := netwData.Containers
	for _, cont := range containers {
		nsvc.CurrentNodeState.Containers[cont.Name].StopCont(container.StopOptions{})
	}
	return nil
}

func (nsvc *NodeService) restoreContainers(netw cluster.OrchNetwork) error {
	for _, cont := range nsvc.CurrentNodeState.Containers {
		if cont.CurrentStatus != cont.DesiredStatus && cont.CurrentStatus == "stopped" {
			cont.StartCont(types.ContainerStartOptions{})
		}
	}
	return nil
}

func (nsvc *NodeService) changeNetworks() bool {
	var err error
	change := false
	for name, netw := range nsvc.DesiredNodeState.Networks {
		currentNetw := nsvc.CurrentNodeState.Networks[name]
		if currentNetw != nil {
			if !(reflect.DeepEqual(netw.NetworkConfig, currentNetw.NetworkConfig)) {
				nsvc.stopContainersOnNetwork(*currentNetw)
				nsvc.deployNetwork(netw)
				change = true
				nsvc.restoreContainers(*currentNetw)
			} else if netw.DesiredStatus != currentNetw.CurrentStatus {
				currentNetw.DesiredStatus = netw.DesiredStatus
				switch currentNetw.DesiredStatus {
				case "created":
					currentNetw.CreateNet(netw.NetworkConfig)
				case "removed":
					//nsvc.stopContainersOnNetwork(*currentNetw)
					currentNetw.RemoveNet()
					//nsvc.restoreContainers(*currentNetw)
				}
				if currentNetw.CurrentStatus == currentNetw.DesiredStatus && err == nil {
					nsvc.clusterChangeLog.Logger.Info("Successful network operation", "name", currentNetw.Name, "status", currentNetw.CurrentStatus)
				} else {
					nsvc.clusterChangeLog.Logger.Error("Failed network operation", "name", currentNetw.Name, "status", currentNetw.CurrentStatus, "error", err)
				}
				change = true
			}
		} else {
			nsvc.deployNetwork(netw)
			change = true
		}
	}
	return change
}
