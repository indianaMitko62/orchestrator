package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/indianaMitko62/orchestrator/src/master"
)

func main() {
	var err error
	//name := "local node"

	clusterState := master.NewClusterState() // for testing yaml

	nodeManager := &master.NodeManager{
		NodeSettings: master.NodeSettings{
			Name:    "Node1",
			Address: "127.0.0.1",
			Containers: []*master.ContainerConfig{
				{
					Name:  "Container1",
					Image: "nginx:latest",
					NetworkConfig: master.ContainerNetworkConfig{
						NetworkID:   "net1",
						IPv4Address: "172.16.0.2",
						DNS:         []string{"8.8.8.8", "8.8.4.4"},
					},
				},
			},
			Networks: []*master.NetworkConfig{
				{
					NetworkID: "net1",
				},
			},
			VolumeConfig: []*master.VolumeConfig{
				{
					VolumeID: "vol1",
				},
			},
		},
		Client: nil,
	}
	nodeManager1 := &master.NodeManager{
		NodeSettings: master.NodeSettings{
			Name:    "Node2",
			Address: "127.0.0.1",
			Containers: []*master.ContainerConfig{
				{
					Name:  "Container2",
					Image: "nginx:latest",
					NetworkConfig: master.ContainerNetworkConfig{
						NetworkID:   "net2",
						IPv4Address: "172.16.0.2",
						DNS:         []string{"8.8.8.8", "8.8.4.4"},
					},
				},
			},
			Networks: []*master.NetworkConfig{
				{
					NetworkID: "net2",
				},
			},
			VolumeConfig: []*master.VolumeConfig{
				{
					VolumeID: "vol2",
				},
			},
		},
		Client: nil,
	}

	clusterState.Nodes[nodeManager.Name] = nodeManager
	clusterState.Nodes[nodeManager1.Name] = nodeManager1
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
