package node

import (
	"reflect"

	"github.com/indianaMitko62/orchestrator/src/cluster"
)

func (nsvc *NodeService) DeployNewNetwork(netw *cluster.OrchNetwork) {
	netw.Cli = nsvc.cli
	_, err := netw.CreateNet(netw.NetworkConfig)
	if err != nil {
		err = nsvc.HandleDuplicateNetworks(netw)
	}
	if netw.DesiredStatus == netw.CurrentStatus && err == nil {
		nsvc.CurrentNodeState.Networks[netw.Name] = netw
		nsvc.clusterChangeLog.Logger.Info("Network successfully created", "name", netw.Name)
	} else {
		nsvc.clusterChangeLog.Logger.Info("Could not create network", "name", netw.Name, "err", err.Error())
	}
}

func (nsvc *NodeService) HandleDuplicateNetworks(newNet *cluster.OrchNetwork) error {
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

func (nsvc *NodeService) changeNetworks() bool {
	var err error
	change := false
	for name, netw := range nsvc.DesiredNodeState.Networks {
		currentNetw := nsvc.CurrentNodeState.Networks[name]
		if currentNetw != nil {
			if !(reflect.DeepEqual(netw.NetworkConfig, currentNetw.NetworkConfig)) {
				nsvc.DeployNewNetwork(netw)
				change = true
			} else if netw.DesiredStatus != currentNetw.CurrentStatus {
				currentNetw.DesiredStatus = netw.DesiredStatus
				switch currentNetw.DesiredStatus {
				case "created":
					currentNetw.CreateNet(netw.NetworkConfig)
				case "removed":
					currentNetw.RemoveNet()
				}
				if currentNetw.CurrentStatus == currentNetw.DesiredStatus && err == nil {
					nsvc.clusterChangeLog.Logger.Info("Successful network operation", "name", currentNetw.Name, "status", currentNetw.CurrentStatus)
				} else {
					nsvc.clusterChangeLog.Logger.Error("Failed network operation", "name", currentNetw.Name, "status", currentNetw.CurrentStatus, "error", err)
				}
				change = true
			}
		} else {
			nsvc.DeployNewNetwork(netw)
			change = true
		}
	}
	return change
}
