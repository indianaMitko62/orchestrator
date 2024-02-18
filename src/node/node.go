package node

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/indianaMitko62/orchestrator/src/cluster"
	"gopkg.in/yaml.v3"
)

/*
TODO: functions managing overall node performance and loading(cpu, memory, disk) and overall node logic
*/

func (nsvc *NodeService) CompareStates() bool { // true for same, false for different, possibly to identify changes
	return true
}

func (nsvc *NodeService) HandleDuplicateContainers(newCont *cluster.OrchContainer) error {
	slog.Info("Trying to remove duplicate container if exists", "name", newCont.ContainerConfig.Hostname)
	newCont.StopCont(container.StopOptions{})
	newCont.RemoveCont(types.ContainerRemoveOptions{})

	slog.Info("Trying to create container again", "name", newCont.ContainerConfig.Hostname)
	_, err := newCont.CreateCont()
	if err != nil {
		slog.Error("Second attempt for container creation failed. Aborting...", "name", newCont.ContainerConfig.Hostname)
		return err
	}
	return nil
}

func (nsvc *NodeService) HandleDuplicateNetworks(newNet *cluster.OrchNetwork) error {
	slog.Info("Trying to remove duplicate network if exists", "name", newNet.Name)
	newNet.RemoveNet()

	slog.Info("Trying to create network again", "name", newNet.Name)
	_, err := newNet.CreateNet(newNet.NetworkConfig)
	if err != nil {
		slog.Error("Second attempt for network creation failed. Aborting...", "name", newNet.Name)
		return err
	}
	return nil
}

func (nsvc *NodeService) HandleDuplicateVolumes(newVol *cluster.OrchVolume) error {
	slog.Info("Trying to remove duplicate volume if exists", "name", newVol.Name)
	newVol.RemoveVol(true)

	slog.Info("Trying to create volume again", "name", newVol.Name)
	_, err := newVol.CreateVol(newVol.Config)
	if err != nil {
		slog.Error("Second attempt for volume creation failed. Aborting...", "name", newVol.Name)
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
			nsvc.ClusterChangeOutcome.Logs[name] = "sucessfully " + img.CurrentStatus
		} else {
			nsvc.ClusterChangeOutcome.Successful = false
			nsvc.ClusterChangeOutcome.Logs[name] = err.Error()
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
			nsvc.ClusterChangeOutcome.Logs[name] = "successfully " + netw.CurrentStatus
		} else {
			nsvc.ClusterChangeOutcome.Successful = false
			nsvc.ClusterChangeOutcome.Logs[name] = err.Error()
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
			nsvc.ClusterChangeOutcome.Logs[name] = "successfully " + vol.CurrentStatus
		} else {
			nsvc.ClusterChangeOutcome.Successful = false
			nsvc.ClusterChangeOutcome.Logs[name] = err.Error()
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
			nsvc.ClusterChangeOutcome.Logs[name] = "successfully " + cont.CurrentStatus
		} else {
			nsvc.ClusterChangeOutcome.Successful = false
			nsvc.ClusterChangeOutcome.Logs[name] = err.Error()
		}
	}
	fmt.Println()
	for name, log := range nsvc.ClusterChangeOutcome.Logs { // for result
		fmt.Println(name, log)
	}
	nsvc.postClusterChangeOutcome(nsvc.MasterAddress + "/clusterState")
	return nil
}

func (nsvc *NodeService) postClusterChangeOutcome(URL string) {

	yamlData, err := yaml.Marshal(nsvc.ClusterChangeOutcome)
	if err != nil {
		slog.Error("Could not marshall Cluster Change Outcome logs to yaml")
	}
	bodyReader := bytes.NewReader(yamlData)
	req, err := http.NewRequest(http.MethodPost, URL, bodyReader)
	if err != nil {
		slog.Error("Could not create POST request", "URL", URL)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Could not send POST request")
	}

	if resp.StatusCode == http.StatusOK {
		slog.Info("Cluster Change Outcome logs send successfully")
	}

}

// func (nsvc *NodeService) findDifferences() error {
// 	slog.Info("finding differences")
// 	for name, cont := range nsvc.DesiredNodeState.Containers {
// 		if nsvc.CurrentNodeState.Containers[name] != nil {
// 			currentcont := nsvc.CurrentNodeState.Containers[name]
// 			cont.Cli = currentcont.Cli
// 			cont.ID = currentcont.ID
// 			cont.CurrentStatus = currentcont.CurrentStatus
// 			cont.DesiredStatus = currentcont.DesiredStatus
// 			if !reflect.DeepEqual(*cont, *currentcont) { // just for setting change. For status change switch case to be implemented
// 				slog.Info("Change in container", "name", name)
// 				currentcont.StopCont(container.StopOptions{})
// 				currentcont.RemoveCont(types.ContainerRemoveOptions{})
// 				cont.CreateCont()
// 				if cont.DesiredStatus == "running" {
// 					cont.StartCont(types.ContainerStartOptions{})
// 				}
// 				nsvc.CurrentNodeState.Containers[name] = cont
// 			} else {
// 				slog.Info("Same container", "name", name)
// 			}
// 		} else {
// 			cont.Cli = nsvc.cli
// 			cont.CreateCont()
// 			if cont.DesiredStatus == "running" {
// 				cont.StartCont(types.ContainerStartOptions{})
// 			}
// 			slog.Info("Container added", "name", name)
// 		}
// 	}
// 	return nil
// }

func (nsvc *NodeService) Node() error {
	nsvc.MasterAddress = "http://localhost:1986" //harcoded for now
	clusterStateURL := nsvc.MasterAddress + "/clusterState"

	for {
		recievedClusterState, err := cluster.GetClusterState(clusterStateURL)
		if err != nil {
			slog.Error("could not get cluster data", "error", err)
		} else {
			nsvc.DesiredNodeState = &recievedClusterState.Nodes[nsvc.Name].NodeState
			if nsvc.CurrentNodeState == nil {
				slog.Info("No current node state")
				err := nsvc.InitCluster()
				if err != nil {
					slog.Error("Could not init cluster")
				}
			}
			slog.Info("Present current node state")
			//go nsvc.findDifferences()
		}
		slog.Info("Main Node process sleeping...")
		time.Sleep(time.Second * 5)
		fmt.Print("\n\n\n")
	}
}
