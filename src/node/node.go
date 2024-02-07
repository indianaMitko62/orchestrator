package node

import (
	"log/slog"
	"strings"

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

func (nsvc *NodeService) InitCluster() error {
	nsvc.CurrentNodeState = cluster.NewNodeState()
	for name, img := range nsvc.DesiredNodeState.Images {
		img.Cli = nsvc.cli
		imgNode := *img
		imgNode.PullImg(&types.ImagePullOptions{ // to work you must read pullImg return value and init cli somehow. line 23 may not be the best way. Cli to be moved ?
			All:           imgNode.All,
			RegistryAuth:  imgNode.RegistryAuth,
			Platform:      imgNode.Platform,
			PrivilegeFunc: nil,
		})
		nsvc.CurrentNodeState.Images[name] = &imgNode
	}
	for name, network := range nsvc.DesiredNodeState.Networks {
		network.Cli = nsvc.cli
		netNode := *network
		netNode.CreateNet(netNode.NetworkConfig)
		nsvc.CurrentNodeState.Networks[name] = &netNode
	}
	for name, vol := range nsvc.DesiredNodeState.Volumes {
		vol.Cli = nsvc.cli
		volNode := *vol
		volNode.CreateVol(volNode.Config)
		nsvc.CurrentNodeState.Volumes[name] = &volNode
	}
	for name, cont := range nsvc.DesiredNodeState.Containers {
		cont.Cli = nsvc.cli
		contNode := *cont
		defer func() {
			if r := recover(); r != nil {
				slog.Info("Trying to remove and stop duplicate container if exists", "name", contNode.ContainerConfig.Hostname)
				contNode.StopCont(container.StopOptions{})
				contNode.RemoveCont(types.ContainerRemoveOptions{})
				slog.Info("Trying to create container again", "name", contNode.ContainerConfig.Hostname)
				_, err := contNode.CreateCont()
				if err != nil {
					slog.Error("Second attempt for container creation failed. Aborting...", "name", contNode.ContainerConfig.Hostname)
					return
				}
				if strings.ToLower(cont.Status) == "running" {
					contNode.StartCont(types.ContainerStartOptions{})
				}
				nsvc.CurrentNodeState.Containers[name] = &contNode
			}
		}()
		_, err := contNode.CreateCont()
		if err != nil {
			panic(err)
		}
		if strings.ToLower(cont.Status) == "running" {
			contNode.StartCont(types.ContainerStartOptions{})
		}
	}
	return nil
}

func (nsvc *NodeService) Node() error {
	slog.Info("Entered")
	nsvc.MasterAddress = "localhost" //harcoded for now
	clusterStateURL := "http://" + nsvc.MasterAddress + ":1986/clusterState"

	recievedClusterState, err := cluster.GetClusterState(clusterStateURL)
	if err != nil {
		slog.Error("could not get cluster data", "error", err)
	} else {
		nsvc.DesiredNodeState = &recievedClusterState.Nodes[nsvc.Name].NodeState
		//nsvc.DesiredNodeState.Cli, _ = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		slog.Info("alo da", "DS", nsvc.DesiredNodeState, "CS", nsvc.CurrentNodeState)
		if nsvc.CurrentNodeState == nil {
			slog.Info("No current node state")
			err := nsvc.InitCluster()
			if err != nil {
				slog.Error("Could not init cluster")
			}
			nsvc.CurrentNodeState = nsvc.DesiredNodeState
		}
	}
	return nil
}
