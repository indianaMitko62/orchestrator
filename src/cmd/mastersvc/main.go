package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

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

	nodeManager := cluster.NodeManager{
		NodeSettings: cluster.NodeSettings{
			Name: "Node1",
		},
		NodeState: cluster.NodeState{
			Containers: map[string]*cluster.OrchContainer{
				"Container1": {
					DesiredStatus: "running",
					ContainerConfig: &container.Config{
						Hostname:     "Container1",
						Image:        "nginx:latest",
						ExposedPorts: map[nat.Port]struct{}{"80/tcp": {}},
						Healthcheck: &container.HealthConfig{
							Test:     []string{"CMD", "echo", "0"}, // vinagi ama vinagi CMD.
							Interval: 5 * time.Second,
							Timeout:  2 * time.Second,
						},
					},
					HostConfig: &container.HostConfig{
						PortBindings: nat.PortMap{
							"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8080"}},
						},
					},
					NetworkingConfig: &network.NetworkingConfig{},
				},
				"Container2": {
					DesiredStatus: "running",
					ContainerConfig: &container.Config{
						Image:        "nginx:latest",
						Hostname:     "Container2",
						ExposedPorts: map[nat.Port]struct{}{"80/tcp": {}},
						Healthcheck: &container.HealthConfig{
							Test:     []string{"CMD", "echo", "0"}, // vinagi ama vinagi CMD.
							Interval: 5 * time.Second,
							Timeout:  2 * time.Second,
						},
					},
					HostConfig: &container.HostConfig{
						PortBindings: nat.PortMap{
							"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8081"}},
						},
					},
					NetworkingConfig: &network.NetworkingConfig{},
				},
				"Container3": {
					DesiredStatus: "stopped",
					ContainerConfig: &container.Config{
						Hostname:     "Container3",
						Image:        "nginx:latest",
						ExposedPorts: map[nat.Port]struct{}{"80/tcp": {}},
						Healthcheck: &container.HealthConfig{
							Test:     []string{"CMD", "echo", "0"}, // vinagi ama vinagi CMD.
							Interval: 5 * time.Second,
							Timeout:  2 * time.Second,
						},
					},
					HostConfig: &container.HostConfig{
						PortBindings: nat.PortMap{
							"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8083"}},
						},
					},
					NetworkingConfig: &network.NetworkingConfig{},
				},
			},
			Networks: map[string]*cluster.OrchNetwork{
				"indiana net": {
					Name:          "indiana net",
					DesiredStatus: "created",
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
	nodeManager1 := cluster.NodeManager{
		NodeSettings: cluster.NodeSettings{
			Name: "Node2",
		},
		NodeState: cluster.NodeState{
			Containers: map[string]*cluster.OrchContainer{
				"Container21": {
					DesiredStatus: "running",
					ContainerConfig: &container.Config{
						Hostname:     "Container21",
						Image:        "nginx:latest",
						ExposedPorts: map[nat.Port]struct{}{"80/tcp": {}},
						Healthcheck: &container.HealthConfig{
							Test:     []string{"CMD", "echo", "0"}, // vinagi ama vinagi CMD.
							Interval: 5 * time.Second,
							Timeout:  2 * time.Second,
						},
					},
					HostConfig: &container.HostConfig{
						PortBindings: nat.PortMap{
							"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8090"}},
						},
					},
					NetworkingConfig: &network.NetworkingConfig{},
				},
				"Container22": {
					DesiredStatus: "running",
					ContainerConfig: &container.Config{
						Image:        "nginx:latest",
						Hostname:     "Container22",
						ExposedPorts: map[nat.Port]struct{}{"80/tcp": {}},
						Healthcheck: &container.HealthConfig{
							Test:     []string{"CMD", "echo", "0"}, // vinagi ama vinagi CMD.
							Interval: 5 * time.Second,
							Timeout:  2 * time.Second,
						},
					},
					HostConfig: &container.HostConfig{
						PortBindings: nat.PortMap{
							"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8091"}},
						},
					},
					NetworkingConfig: &network.NetworkingConfig{},
				},
				"Container23": {
					DesiredStatus: "stopped",
					ContainerConfig: &container.Config{
						Hostname:     "Container23",
						Image:        "nginx:latest",
						ExposedPorts: map[nat.Port]struct{}{"80/tcp": {}},
						Healthcheck: &container.HealthConfig{
							Test:     []string{"CMD", "echo", "0"}, // vinagi ama vinagi CMD.
							Interval: 5 * time.Second,
							Timeout:  2 * time.Second,
						},
					},
					HostConfig: &container.HostConfig{
						PortBindings: nat.PortMap{
							"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8093"}},
						},
					},
					NetworkingConfig: &network.NetworkingConfig{},
				},
			},
			Networks: map[string]*cluster.OrchNetwork{
				"indiana net2": {
					Name:          "indiana net2",
					DesiredStatus: "created",
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

	clusterState.Nodes[nodeManager.Name] = nodeManager
	clusterState.Nodes[nodeManager1.Name] = nodeManager1
	clusterState.CollectImages() // to be developed and added to master logic
	yamlData, _ := cluster.ToYaml(clusterState)
	fmt.Printf("%s", yamlData)

	msvc := master.NewMasterService()
	msvc.CS = clusterState
	//err = msvc.ConnectToNodes()
	if err != nil {
		slog.Error("could not connect to nodes", "err", err)
		os.Exit(1)
	}

	msvc.Master()
}
