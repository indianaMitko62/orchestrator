package node

import (
	"reflect"

	"github.com/indianaMitko62/orchestrator/src/cluster"
)

func (nsvc *NodeService) handleDuplicateVolumes(newVol *cluster.OrchVolume) error {
	nsvc.nodeLog.Logger.Info("Trying to remove duplicate volume if exists", "name", newVol.Name)
	newVol.RemoveVol(true)

	nsvc.nodeLog.Logger.Info("Trying to create volume again", "name", newVol.Name)
	_, err := newVol.CreateVol(newVol.Config)
	if err != nil {
		nsvc.nodeLog.Logger.Error("Second attempt for volume creation failed. Aborting...", "name", newVol.Name)
		return err
	}
	return nil
}

func (nsvc *NodeService) deployVolume(vol *cluster.OrchVolume) {
	vol.Cli = nsvc.cli
	_, err := vol.CreateVol(vol.Config)
	if err != nil {
		err = nsvc.handleDuplicateVolumes(vol)
	}
	if vol.DesiredStatus == vol.CurrentStatus && err == nil {
		nsvc.CurrentNodeState.Volumes[vol.Name] = vol
		nsvc.clusterChangeLog.Logger.Info("Volume successfully created", "name", vol.Name)
	} else {
		nsvc.clusterChangeLog.Logger.Info("Could not create volume", "name", vol.Name, "err", err.Error())
	}
}

func (nsvc *NodeService) changeVolumes() bool {
	var err error
	change := false
	for name, vol := range nsvc.DesiredNodeState.Volumes {
		currentVol := nsvc.CurrentNodeState.Volumes[name]
		if currentVol != nil {
			if !(reflect.DeepEqual(vol.Config, currentVol.Config)) {
				nsvc.deployVolume(vol)
				change = true
			} else if vol.DesiredStatus != currentVol.CurrentStatus {
				currentVol.DesiredStatus = vol.DesiredStatus
				switch currentVol.DesiredStatus {
				case "created":
					_, err = currentVol.CreateVol(vol.Config)
				case "removed":
					err = currentVol.RemoveVol(true)
				}
				if currentVol.CurrentStatus == currentVol.DesiredStatus && err == nil {
					nsvc.clusterChangeLog.Logger.Info("Successful volume operation", "name", currentVol.Name, "status", currentVol.CurrentStatus)
				} else {
					nsvc.clusterChangeLog.Logger.Error("Failed volume operation", "name", currentVol.Name, "status", currentVol.CurrentStatus, "error", err)
				}
				change = true
			}
		} else {
			nsvc.deployVolume(vol)
			change = true
		}
	}
	return change
}
