package node

import (
	"reflect"

	"github.com/docker/docker/api/types"
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

func (nsvc *NodeService) changeNetworks() bool {
	var err error
	change := false
	for name, netw := range nsvc.DesiredNodeState.Networks {
		currentNetw := nsvc.CurrentNodeState.Networks[name]
		if currentNetw != nil {
			if !(reflect.DeepEqual(netw.NetworkConfig, currentNetw.NetworkConfig)) {
				netwData, _ := currentNetw.InspectNet(types.NetworkInspectOptions{})
				if err != nil {
					nsvc.nodeLog.Logger.Error("Could not check network for active endpoints")
					continue
				}
				if len(netwData.Containers) > 0 {
					nsvc.clusterChangeLog.Logger.Error("Cannot change network that has active endpoints")
					continue
				}
				nsvc.deployNetwork(netw)
				change = true
			} else if netw.DesiredStatus != currentNetw.CurrentStatus {
				currentNetw.DesiredStatus = netw.DesiredStatus
				switch currentNetw.DesiredStatus {
				case "created":
					_, err = currentNetw.CreateNet(netw.NetworkConfig)
				case "removed":
					err = currentNetw.RemoveNet()
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
