package node

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

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

func (nsvc *NodeService) HandleDuplicateContainers(contNode cluster.OrchContainer, cont cluster.OrchContainer, name string) error {
	slog.Info("Trying to stop and remove duplicate container if exists", "name", cont.ContainerConfig.Hostname)
	cont.StopCont(container.StopOptions{})
	cont.RemoveCont(types.ContainerRemoveOptions{})
	slog.Info("Trying to create container again", "name", cont.ContainerConfig.Hostname)
	_, err := contNode.CreateCont()
	if err != nil {
		slog.Error("Second attempt for container creation failed. Aborting...", "name", contNode.ContainerConfig.Hostname)
		return err
	}
	return nil
}

func (nsvc *NodeService) HandleDuplicateNetworks(netNode cluster.OrchNetwork, net cluster.OrchNetwork) error {
	slog.Info("Trying to remove duplicate network if exists", "name", net.Name)
	net.RemoveNet()
	slog.Info("Trying to create network again", "name", net.Name)
	_, err := netNode.CreateNet(netNode.NetworkConfig)
	if err != nil {
		slog.Error("Second attempt for network creation failed. Aborting...", "name", netNode.Name)
		return err
	}
	return nil
}

func (nsvc *NodeService) HandleDuplicateVolumes(volNode cluster.OrchVolume, vol cluster.OrchVolume) error {
	slog.Info("Trying to remove duplicate volume if exists", "name", vol.Name)
	vol.RemoveVol(true)
	slog.Info("Trying to create volume again", "name", vol.Name)
	_, err := volNode.CreateVol(volNode.Config)
	if err != nil {
		slog.Error("Second attempt for volume creation failed. Aborting...", "name", volNode.Name)
		return err
	}
	return nil
}

func (nsvc *NodeService) InitCluster() error {
	nsvc.CurrentNodeState = cluster.NewNodeState()
	fmt.Println()

	for name, img := range nsvc.DesiredNodeState.Images {
		img.Cli = nsvc.cli
		imgNode := *img
		_, err := imgNode.PullImg(&types.ImagePullOptions{
			All:           imgNode.All,
			RegistryAuth:  imgNode.RegistryAuth,
			Platform:      imgNode.Platform,
			PrivilegeFunc: nil,
		})
		if err == nil {
			nsvc.CurrentNodeState.Images[name] = &imgNode
			nsvc.ClusterChangeOutcome.Logs[name] = "sucessfully " + imgNode.Status
			continue
		}
		nsvc.ClusterChangeOutcome.Successful = false
		nsvc.ClusterChangeOutcome.Logs[name] = err.Error()
	}
	fmt.Println()

	for name, network := range nsvc.DesiredNodeState.Networks {
		network.Cli = nsvc.cli
		netNode := *network
		_, err := netNode.CreateNet(netNode.NetworkConfig)
		if err != nil {
			err = nsvc.HandleDuplicateNetworks(netNode, *network)
		}
		if netNode.Status == network.Status && err == nil {
			nsvc.CurrentNodeState.Networks[name] = &netNode
			nsvc.ClusterChangeOutcome.Logs[name] = "successfully " + netNode.Status
			continue
		}
		nsvc.ClusterChangeOutcome.Successful = false
		nsvc.ClusterChangeOutcome.Logs[name] = err.Error()
	}
	fmt.Println()

	for name, vol := range nsvc.DesiredNodeState.Volumes {
		vol.Cli = nsvc.cli
		volNode := *vol
		_, err := volNode.CreateVol(volNode.Config)
		if err != nil {
			err = nsvc.HandleDuplicateVolumes(volNode, *vol)
		}
		if volNode.Status == vol.Status && err == nil {
			nsvc.CurrentNodeState.Volumes[name] = &volNode
			nsvc.ClusterChangeOutcome.Logs[name] = "successfully " + volNode.Status
			continue
		}
		nsvc.ClusterChangeOutcome.Successful = false
		nsvc.ClusterChangeOutcome.Logs[name] = err.Error()
	}
	fmt.Println()

	for name, cont := range nsvc.DesiredNodeState.Containers {
		cont.Cli = nsvc.cli
		contNode := *cont

		_, err := contNode.CreateCont()
		if err != nil {
			nsvc.HandleDuplicateContainers(contNode, *cont, name)
		}
		if strings.ToLower(cont.Status) == "running" {
			contNode.StartCont(types.ContainerStartOptions{})
		}
		if contNode.Status == cont.Status {
			nsvc.CurrentNodeState.Containers[name] = &contNode
			nsvc.ClusterChangeOutcome.Logs[name] = "successfully " + contNode.Status
			continue
		}
		nsvc.ClusterChangeOutcome.Successful = false
		nsvc.ClusterChangeOutcome.Logs[name] = err.Error()
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

func (nsvc *NodeService) Node() error {
	nsvc.MasterAddress = "http://localhost:1986" //harcoded for now
	clusterStateURL := nsvc.MasterAddress + "/clusterState"

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
	}
	return nil
}
