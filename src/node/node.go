package node

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

/*
TODO: functions managing overall node performance and loading(cpu, memory, disk) and overall node logic
*/

func (nsvc *NodeService) CompareStates() bool { // true for same, false for different, possibly to identify changes
	return true
}

func (nsvc *NodeService) HandleDuplicateContainers(newCont *cluster.OrchContainer) error {
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

func (nsvc *NodeService) HandleDuplicateVolumes(newVol *cluster.OrchVolume) error {
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

func (nsvc *NodeService) DeployNewContainer(cont *cluster.OrchContainer) {
	cont.Cli = nsvc.cli
	_, err := cont.CreateCont()
	if err != nil {
		err = nsvc.HandleDuplicateContainers(cont)
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

func (nsvc *NodeService) DeployNewVolume(vol *cluster.OrchVolume) {
	vol.Cli = nsvc.cli
	_, err := vol.CreateVol(vol.Config)
	if err != nil {
		err = nsvc.HandleDuplicateVolumes(vol)
	}
	if vol.DesiredStatus == vol.CurrentStatus && err == nil {
		nsvc.CurrentNodeState.Volumes[vol.Name] = vol
		nsvc.clusterChangeLog.Logger.Info("Volume successfully created", "name", vol.Name)
	} else {
		nsvc.clusterChangeLog.Logger.Info("Could not create volume", "name", vol.Name, "err", err.Error())
	}
}

func (nsvc *NodeService) InitCluster() error { // status change. Refactor function needed
	nsvc.CurrentNodeState = cluster.NewNodeState()
	fmt.Println()

	for name, img := range nsvc.DesiredNodeState.Images {
		img.Cli = nsvc.cli
		_, err := img.PullImg(&types.ImagePullOptions{
			All:           img.All,
			RegistryAuth:  img.RegistryAuth,
			Platform:      img.Platform,
			PrivilegeFunc: nil,
		})
		if img.CurrentStatus == img.DesiredStatus && err == nil {
			nsvc.CurrentNodeState.Images[name] = img
			nsvc.clusterChangeLog.Logger.Info("Image successfully pulled", "name", img.Name)
		} else {
			nsvc.clusterChangeLog.Logger.Info("Could not pull image", "name", img.Name, "err", err.Error())
		}
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

func (nsvc *NodeService) findDifferences() error { // add error handling and logging
	nsvc.nodeLog.Logger.Info("finding differences")
	change := false
	for name, cont := range nsvc.DesiredNodeState.Containers {
		currentCont := nsvc.CurrentNodeState.Containers[name]
		if currentCont != nil {
			if !(reflect.DeepEqual(cont.ContainerConfig, currentCont.ContainerConfig) &&
				reflect.DeepEqual(cont.HostConfig, currentCont.HostConfig) && reflect.DeepEqual(cont.NetworkingConfig, currentCont.NetworkingConfig)) {
				nsvc.DeployNewContainer(cont)
				change = true
			} else if cont.DesiredStatus != currentCont.CurrentStatus {
				currentCont.DesiredStatus = cont.DesiredStatus
				var err error
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
			nsvc.DeployNewContainer(cont)
			change = true
		}
	}
	if change {
		nsvc.postClusterChangeOutcome(nsvc.MasterAddress + "/clusterState")
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
