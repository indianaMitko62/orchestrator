package node

import (
	"fmt"
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
	nsvc.clusterChangeLog.Info("Trying to remove duplicate container if exists", "name", newCont.ContainerConfig.Hostname)
	newCont.StopCont(container.StopOptions{})
	newCont.RemoveCont(types.ContainerRemoveOptions{})

	nsvc.clusterChangeLog.Info("Trying to create container again", "name", newCont.ContainerConfig.Hostname)
	_, err := newCont.CreateCont()
	if err != nil {
		nsvc.clusterChangeLog.Error("Second attempt for container creation failed. Aborting...", "name", newCont.ContainerConfig.Hostname)
		return err
	}
	return nil
}

func (nsvc *NodeService) HandleDuplicateNetworks(newNet *cluster.OrchNetwork) error {
	nsvc.clusterChangeLog.Info("Trying to remove duplicate network if exists", "name", newNet.Name)
	newNet.RemoveNet()

	nsvc.clusterChangeLog.Info("Trying to create network again", "name", newNet.Name)
	_, err := newNet.CreateNet(newNet.NetworkConfig)
	if err != nil {
		nsvc.clusterChangeLog.Error("Second attempt for network creation failed. Aborting...", "name", newNet.Name)
		return err
	}
	return nil
}

func (nsvc *NodeService) HandleDuplicateVolumes(newVol *cluster.OrchVolume) error {
	nsvc.clusterChangeLog.Info("Trying to remove duplicate volume if exists", "name", newVol.Name)
	newVol.RemoveVol(true)

	nsvc.clusterChangeLog.Info("Trying to create volume again", "name", newVol.Name)
	_, err := newVol.CreateVol(newVol.Config)
	if err != nil {
		nsvc.clusterChangeLog.Error("Second attempt for volume creation failed. Aborting...", "name", newVol.Name)
		return err
	}
	return nil
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
			nsvc.clusterChangeLog.Info("Image successfully pulled", "name", img.Name)
		} else {
			nsvc.clusterChangeLog.Info("Could not pull image", "name", img.Name, "err", err.Error())
		}
	}
	fmt.Println()

	for name, netw := range nsvc.DesiredNodeState.Networks {
		netw.Cli = nsvc.cli
		_, err := netw.CreateNet(netw.NetworkConfig)
		if err != nil {
			err = nsvc.HandleDuplicateNetworks(netw)
		}
		if netw.DesiredStatus == netw.CurrentStatus && err == nil {
			nsvc.CurrentNodeState.Networks[name] = netw
			nsvc.clusterChangeLog.Info("Network successfully created", "name", netw.Name)
		} else {
			nsvc.clusterChangeLog.Info("Could not create network", "name", netw.Name, "err", err.Error())
		}
	}
	fmt.Println()

	for name, vol := range nsvc.DesiredNodeState.Volumes {
		vol.Cli = nsvc.cli
		_, err := vol.CreateVol(vol.Config)
		if err != nil {
			err = nsvc.HandleDuplicateVolumes(vol)
		}
		if vol.DesiredStatus == vol.CurrentStatus && err == nil {
			nsvc.CurrentNodeState.Volumes[name] = vol
			nsvc.clusterChangeLog.Info("Volume successfully created", "name", vol.Name)
		} else {
			nsvc.clusterChangeLog.Info("Could not create volume", "name", vol.Name, "err", err.Error())

		}
	}
	fmt.Println()

	for name, cont := range nsvc.DesiredNodeState.Containers {
		cont.Cli = nsvc.cli
		_, err := cont.CreateCont()
		if err != nil {
			err = nsvc.HandleDuplicateContainers(cont)
		}
		if cont.DesiredStatus == "running" {
			err = cont.StartCont(types.ContainerStartOptions{})
		}
		if cont.DesiredStatus == cont.CurrentStatus && err == nil {
			nsvc.CurrentNodeState.Containers[name] = cont
			nsvc.clusterChangeLog.Info("Container successfully created", "name", cont.ContainerConfig.Hostname)
		} else {
			nsvc.clusterChangeLog.Info("Could not create container", "name", cont.ContainerConfig.Hostname, "err", err.Error())

		}
	}
	fmt.Println()
	nsvc.postClusterChangeOutcome(nsvc.MasterAddress + "/clusterState")
	return nil
}

func (nsvc *NodeService) postClusterChangeOutcome(URL string) {

	bodyReader := nsvc.logReader
	req, err := http.NewRequest(http.MethodPost, URL, bodyReader)
	if err != nil {
		nsvc.clusterChangeLog.Error("Could not create POST request", "URL", URL)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		nsvc.clusterChangeLog.Error("Could not send POST request")
	}

	if resp.StatusCode == http.StatusOK {
		nsvc.clusterChangeLog.Info("Cluster Change Outcome logs send successfully")
	}
	nsvc.logReader, _ = os.Open("./dat2")
}

func (nsvc *NodeService) findDifferences() error { // add error handling and logging
	nsvc.clusterChangeLog.Info("finding differences")
	change := false
	for name, cont := range nsvc.DesiredNodeState.Containers {
		cont.Cli = nsvc.cli
		currentCont := nsvc.CurrentNodeState.Containers[name]
		if currentCont != nil {
			if !(reflect.DeepEqual(cont.ContainerConfig, currentCont.ContainerConfig) &&
				reflect.DeepEqual(cont.HostConfig, currentCont.HostConfig) &&
				reflect.DeepEqual(cont.NetworkingConfig, currentCont.NetworkingConfig)) {

				nsvc.clusterChangeLog.Info("Container settings changed", "name", cont.ContainerConfig.Hostname)
				currentCont.StopCont(container.StopOptions{})
				currentCont.RemoveCont(types.ContainerRemoveOptions{})
				cont.CreateCont()
				if cont.DesiredStatus == "running" {
					cont.StartCont(types.ContainerStartOptions{})
				}
				if cont.DesiredStatus == cont.CurrentStatus {
					nsvc.CurrentNodeState.Containers[name] = cont
				}
				change = true
			}
			//if none check for desired status change - act accordingly
		} else {
			_, err := cont.CreateCont()
			if err != nil {
				err = nsvc.HandleDuplicateContainers(cont)
			}
			if cont.DesiredStatus == "running" && err == nil {
				cont.StartCont(types.ContainerStartOptions{})
			}
			if cont.DesiredStatus == cont.CurrentStatus {
				nsvc.CurrentNodeState.Containers[name] = cont
			}
			nsvc.clusterChangeLog.Info("Container added", "name", name)
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
			nsvc.clusterChangeLog.Error("could not get cluster data", "error", err)
		} else {
			nsvc.DesiredNodeState = &recievedClusterState.Nodes[nsvc.Name].NodeState
			if nsvc.CurrentNodeState == nil {
				nsvc.clusterChangeLog.Info("No current node state")
				err := nsvc.InitCluster()
				if err != nil {
					nsvc.clusterChangeLog.Error("Could not init cluster")
				}
			}
			nsvc.clusterChangeLog.Info("Present current node state")
			go nsvc.findDifferences()
		}
		nsvc.clusterChangeLog.Info("Main Node process sleeping...") // to be moved to different logger, not the one send to /clusterState
		time.Sleep(time.Second * 5)
		fmt.Print("\n\n\n")
	}
}
