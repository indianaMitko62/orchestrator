package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func runCont(cli client.Client) string {

	ctx := context.Background()
	var err error
	resp1, err := cli.ImagePull(ctx, "alpine", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	read, err := io.ReadAll(resp1)
	if err != nil {
		panic(err)
	}

	fmt.Print(string(read))

	containerConfig := &container.Config{
		Image: "alpine:latest",
		Cmd:   []string{"ping", "localhost"},
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}

func listCont(cli client.Client) {

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Container list: ")
	if len(containers) > 0 {
		for _, container := range containers {
			fmt.Printf("%s %s\n", container.ID[:10], container.Image)
		}
	} else {
		fmt.Println("No containers running")
	}
}

func main() {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)

	}
	go listCont(*cli)
	contID := runCont(*cli)

	out, err := cli.ContainerLogs(context.Background(), contID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
