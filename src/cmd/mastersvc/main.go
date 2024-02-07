package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
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
		NodeState: cluster.NodeState{
			Containers: map[string]*cluster.OrchContainer{
				"Container1": {
					Status: "running",
					ContainerConfig: &container.Config{
						Hostname:     "Container1",
						Image:        "nginx:latest",
						ExposedPorts: map[nat.Port]struct{}{"80/tcp": {}},
					},
					HostConfig: &container.HostConfig{
						PortBindings: nat.PortMap{
							"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8080"}},
						},
					},
					NetworkingConfig: &network.NetworkingConfig{},
				},
			},
			Networks: map[string]*cluster.OrchNetwork{
				"net1": {
					ID:     "net ID 1",
					Name:   "indiana net",
					Status: "created",
					NetworkConfig: types.NetworkCreate{
						Driver:         "bridge",
						CheckDuplicate: true,
					},
				},
			},
			Volumes: map[string]*cluster.OrchVolume{},
			Images:  map[string]*cluster.OrchImage{},
		},
	}
	nodeManager1 := &cluster.NodeManager{
		NodeSettings: cluster.NodeSettings{
			Name:    "Node2",
			Address: "127.0.0.1",
		},
		NodeState: cluster.NodeState{
			Containers: map[string]*cluster.OrchContainer{
				"Container2": {
					ContainerConfig: &container.Config{
						Image:        "nginx:latest",
						Hostname:     "Container2",
						ExposedPorts: map[nat.Port]struct{}{"80/tcp": {}},
					},
					HostConfig: &container.HostConfig{
						PortBindings: nat.PortMap{
							"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8080"}},
						},
					},
					NetworkingConfig: &network.NetworkingConfig{},
				},
			},
			Networks: map[string]*cluster.OrchNetwork{
				"net2": {
					ID: "net ID 2",
				},
			},
			Volumes: map[string]*cluster.OrchVolume{
				"vol2": {
					Name: "vol2",
				},
			},
			Images: map[string]*cluster.OrchImage{},
		},
	}

	clusterState.Nodes[nodeManager.Name] = nodeManager
	clusterState.Nodes[nodeManager1.Name] = nodeManager1
	clusterState.CollectImages() // to be developed and added to master logic
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
