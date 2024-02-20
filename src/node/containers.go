package node

import (
	"fmt"
	"reflect"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

func (nsvc *NodeService) handleDuplicateContainers(newCont *cluster.OrchContainer) error {
	nsvc.nodeLog.Logger.Info("Trying to remove duplicate container if exists", "name", newCont.ContainerConfig.Hostname)
	newCont.StopCont(container.StopOptions{})
	newCont.RemoveCont(types.ContainerRemoveOptions{})

	nsvc.nodeLog.Logger.Info("Trying to create container again", "name", newCont.ContainerConfig.Hostname)
	_, err := newCont.CreateCont()
	if err != nil {
		nsvc.nodeLog.Logger.Error("Second attempt for container creation failed. Aborting...", "name", newCont.ContainerConfig.Hostname)
		return err
	}
	return nil
}

func (nsvc *NodeService) deployNewContainer(cont *cluster.OrchContainer) {
	cont.Cli = nsvc.cli
	_, err := cont.CreateCont()
	if err != nil {
		err = nsvc.handleDuplicateContainers(cont)
	}
	if cont.DesiredStatus == "running" && err == nil {
		err = cont.StartCont(types.ContainerStartOptions{})
	}
	if cont.DesiredStatus == cont.CurrentStatus && err == nil {
		nsvc.CurrentNodeState.Containers[cont.ContainerConfig.Hostname] = cont
		nsvc.clusterChangeLog.Logger.Info("Container successfully "+cont.CurrentStatus, "name", cont.ContainerConfig.Hostname, "status", cont.CurrentStatus)
	} else {
		nsvc.clusterChangeLog.Logger.Info("Could not create container", "name", cont.ContainerConfig.Hostname)
	}
}

func (nsvc *NodeService) inspectContainer(cont *cluster.OrchContainer) {
	contInfo, err := cont.InspectCont()
	fmt.Printf("baba: %+v\n", contInfo.ContainerJSONBase.State)
	health := contInfo.ContainerJSONBase.State.Health.Status
	if err == nil {
		nsvc.nodeLog.Logger.Info("Container health: ", "name", cont.ContainerConfig.Hostname, "status", health)
		switch health {
		case "none", "starting":
			return
		case "healthy":
			// add to node status
		case "unhealthy":
			// restart container. just use deployNewContainer method??
			// add to node status
		}
	}
}

func (nsvc *NodeService) changeContainers() bool {
	var err error
	change := false
	for name, cont := range nsvc.DesiredNodeState.Containers {
		currentCont := nsvc.CurrentNodeState.Containers[name]
		if currentCont != nil {
			if !(reflect.DeepEqual(cont.ContainerConfig, currentCont.ContainerConfig) &&
				reflect.DeepEqual(cont.HostConfig, currentCont.HostConfig) && reflect.DeepEqual(cont.NetworkingConfig, currentCont.NetworkingConfig)) {
				nsvc.deployNewContainer(cont)
				change = true
			} else if cont.DesiredStatus != currentCont.CurrentStatus {
				currentCont.DesiredStatus = cont.DesiredStatus
				switch currentCont.DesiredStatus {
				case "running":
					err = currentCont.StartCont(types.ContainerStartOptions{})
					if err != nil {
						currentCont.CreateCont()
						err = currentCont.StartCont(types.ContainerStartOptions{})
					}
				case "stopped":
					err = currentCont.StopCont(container.StopOptions{})
				case "killed":
					err = currentCont.KillCont("")
				case "paused":
					err = currentCont.PauseCont()
				case "unpause":
					err = currentCont.UnpauseCont()
				case "removed":
					if currentCont.CurrentStatus == "running" {
						err = currentCont.StopCont(container.StopOptions{})
					}
					if err == nil {
						err = currentCont.RemoveCont(types.ContainerRemoveOptions{})
					}
				}
				if currentCont.CurrentStatus == currentCont.DesiredStatus && err == nil {
					nsvc.clusterChangeLog.Logger.Info("Successful container operation", "name", currentCont.ContainerConfig.Hostname, "status", currentCont.CurrentStatus)
				} else {
					nsvc.clusterChangeLog.Logger.Error("Failed container operation", "name", currentCont.ContainerConfig.Hostname, "status", currentCont.CurrentStatus, "error", err)
				}
				change = true
			}
		} else {
			nsvc.deployNewContainer(cont)
			change = true
		}
	}
	return change
}
