package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/indianaMitko62/orchestrator/src/cluster"
	"github.com/indianaMitko62/orchestrator/src/master"
)

func main() {
	var err error
	//name := "local node"

	clusterState := cluster.NewClusterState() // for testing yaml

	nodeManager := &cluster.NodeManager{
		NodeSettings: cluster.NodeSettings{
			Name:    "Node1",
			Address: "127.0.0.1",
		},
		Containers: map[string]*cluster.ContainerConfig{
			"Container1": {
				Name:       "Container1",
				Image:      "hello-world",
				Privileged: true,
				NetworkConfig: map[string]*cluster.ContainerNetworkConfig{
					"bridge": {
						NetworkID:   "bridge",
						IPv4Address: "172.16.0.2",
					},
				},
			},
		},
		Networks: map[string]*cluster.NetworkConfig{
			"net1": {
				NetworkID: "net ID 1",
			},
		},
		Volumes: map[string]*cluster.VolumeConfig{
			"vol1": {
				VolumeID: "vol ID 1",
			},
		},
		Images: map[string]*cluster.ImageConfig{},
		Client: nil,
	}
	nodeManager1 := &cluster.NodeManager{
		NodeSettings: cluster.NodeSettings{
			Name:    "Node2",
			Address: "127.0.0.1",
		},
		Containers: map[string]*cluster.ContainerConfig{
			"Container2": {
				Name:  "Container2",
				Image: "nginx:lastest",
				NetworkConfig: map[string]*cluster.ContainerNetworkConfig{
					"bridge": {
						NetworkID:   "bridge",
						IPv4Address: "172.16.0.2",
					},
				},
			},
		},
		Networks: map[string]*cluster.NetworkConfig{
			"net2": {
				NetworkID: "vol ID 2",
			},
		},
		Volumes: map[string]*cluster.VolumeConfig{
			"vol2": {
				VolumeID: "vol ID 2",
			},
		},
		Images: map[string]*cluster.ImageConfig{},
		Client: nil,
	}
	clusterState.Nodes[nodeManager.Name] = nodeManager
	clusterState.Nodes[nodeManager1.Name] = nodeManager1
	clusterState.CollectImages()
	yamlData, _ := clusterState.ToYaml()
	fmt.Printf("%s", yamlData)

	msvc := master.NewMasterService(clusterState)
	//err = msvc.ConnectToNodes()
	if err != nil {
		slog.Error("could not connect to nodes", "err", err)
		os.Exit(1)
	}

	msvc.Master()
}
