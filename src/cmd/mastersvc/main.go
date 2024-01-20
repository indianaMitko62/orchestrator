package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/indianaMitko62/orchestrator/src/master"
	"github.com/indianaMitko62/orchestrator/src/node"
)

func main() {
	var err error
	name := "local node"
	msvc := master.NewMasterService([]*master.NodeSettings{
		{
			Name:    name,
			Address: "localhost:1234",
		},
	})

	err = msvc.ConnectToNodes()
	if err != nil {
		slog.Error("could not connect to nodes", "err", err)
		os.Exit(1)
	}

	//var cont = new(node.Container)
	cont := node.OrchContainer{
		ContainerConfig: &container.Config{
			Image: "alpine:latest",
			Cmd:   []string{"ping", "8.8.8.8"},
		},
		Image_name:       "alpine",
		HostConfig:       nil,
		NetworkingConfig: nil,
		Name:             "cont1", ContID: new(string), ContStatus: new(string)}
	// cont.ContainerConfig = &container.Config{
	// 	Image: "alpine:latest",
	// 	Cmd:   []string{"ping", "8.8.8.8"},
	// }
	// cont.Image_name = "alpine:latest"
	// cont.HostConfig = nil
	// cont.NetworkingConfig = nil
	// cont.ContainerName = "cont1"
	*cont.ContID = "not set"
	*cont.ContStatus = "not set"

	err = msvc.CreateContOn(msvc.Nodes[name], &cont)
	if err != nil {
		slog.Error("could not create container on node", "name", cont.Name, "node", name)
		os.Exit(1)
	}

	err = msvc.StartContOn(msvc.Nodes[name], &cont)
	if err != nil {
		slog.Error("could not start container on node", "name", cont.Name, "node", name)
		os.Exit(1)
	}

	sleep := 20
	slog.Info("Master sleeping", "node_name", name, "duration", sleep)
	time.Sleep(time.Duration(sleep) * time.Second)

	err = msvc.StopContOn(msvc.Nodes[name], &cont)
	if err != nil {
		slog.Error("could not stop container on node", "name", cont.Name, "node", name)
		os.Exit(1)
	}
}
