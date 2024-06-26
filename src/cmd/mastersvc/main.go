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
	"gopkg.in/yaml.v3"
)

func main() {
	var err error
	//name := "local node"

	clusterState := cluster.NewClusterState() // for testing yaml

	nodeState := cluster.NodeState{
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
				NetworkingConfig: &network.NetworkingConfig{
					EndpointsConfig: map[string]*network.EndpointSettings{"indiana_net": {NetworkID: "indiana_net"}},
				},
			},
			"Container3": {
				DesiredStatus: "running",
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
			},
		},
		Networks: map[string]*cluster.OrchNetwork{
			"indiana_net": {
				Name:          "indiana_net",
				DesiredStatus: "created",
				NetworkConfig: types.NetworkCreate{
					Driver:         "bridge",
					CheckDuplicate: true,
					Internal:       false,
				},
			},
		},
		Volumes: map[string]*cluster.OrchVolume{},
		Images:  map[string]*cluster.OrchImage{},
	}
	nodeState1 := cluster.NodeState{
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
			"indiana_net2": {
				Name:          "indiana_net2",
				DesiredStatus: "created",
				NetworkConfig: types.NetworkCreate{
					Driver:         "bridge",
					CheckDuplicate: true,
				},
			},
		},
		Volumes: map[string]*cluster.OrchVolume{},
		Images:  map[string]*cluster.OrchImage{},
	}

	clusterState.Nodes["Node1"] = nodeState
	clusterState.Nodes["Node2"] = nodeState1
	clusterState.CollectImages() // to be developed and added to master logic
	//yamlData, _ := cluster.ToYaml(clusterState)
	//fmt.Printf("%s", yamlData)

	if len(os.Args) < 2 {
		slog.Error("No configuration file passes as command line argument")
	}
	confFile := os.Args[1]
	if confFile == "" {
		slog.Error("No command line argument")
	}
	f, err := os.Open(confFile)
	if err != nil {
		slog.Error("Could not open config file", "name", confFile)
	}

	var masterSet master.MasterSettings
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&masterSet)
	if err != nil {
		slog.Error("Could not decode config file", "name", confFile)
	}
	f.Close()
	fmt.Println(masterSet.DefaultNetworking)
	msvc := master.NewMasterService(masterSet)
	// msvc.CS = clusterState

	msvc.Master()
}
