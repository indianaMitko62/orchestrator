package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var cli *client.Client

type ContainerArgs struct {
	Opts any
}

type Container struct {
	containerConfig  container.Config
	image_name       string
	hostConfig       *container.HostConfig
	networkingConfig *network.NetworkingConfig
	containerName    string
	contID           string
	Opts             any
}

var cont = new(Container)

func (cont *Container) CreateCont(args Container, reply *string) error {

	ctx := context.Background()
	// var err error
	// resp1, err := cli.ImagePull(ctx, cont.image_name, args.Opts.(types.ImagePullOptions))
	// if err != nil {
	// 	return err
	// }

	// read, err := io.ReadAll(resp1)
	// if err != nil {
	// 	return err
	// }

	// fmt.Print(string(read))

	// containerConfig := &container.Config{
	// 	Image: "alpine:latest",
	// 	Cmd:   []string{"ping", "localhost"},
	// }
	fmt.Println(args.Opts)
	resp, err := cli.ContainerCreate(ctx, &args.containerConfig, args.hostConfig, args.networkingConfig, nil, args.containerName)
	if err != nil {
		return err
	}
	cont.contID = resp.ID

	return nil
}

// func (cont *Container) StartCont(args ContainerArgs, reply *string) error {
// 	if err := cli.ContainerStart(context.Background(), cont.contID, args.Opts.(types.ContainerStartOptions)); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (cont *Container) StopCont(args ContainerArgs, reply *string) error {
// 	if err := cli.ContainerStop(context.Background(), cont.containerName, args.Opts.(container.StopOptions)); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func listCont(reply *string) error {

// 	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("Container list: ")
// 	if len(containers) > 0 {
// 		for _, container := range containers {
// 			fmt.Printf("%s %s\n", container.ID[:10], container.Image)
// 		}
// 	} else {
// 		fmt.Println("No containers running")
// 	}
// 	return nil
// }

// func (cont *Container) NewClient(args ContainerArgs, reply *string) error {
// 	var err error
// 	cli, err = client.NewClientWithOpts()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func main() {
	// cont := Container{
	// 	container.Config{
	// 		Image: "alpine:latest",
	// 		Cmd:   []string{"echo", "maikatibe"},
	// 	},
	// 	"alpine", nil, nil, "cont1", ""}
	//go listCont(*cli)
	// cli, err := newClient()
	// if err != nil {
	// 	return err
	// }
	//	contID := cont1.runCont(*cli)
	// cont.createCont(types.ImagePullOptions{})
	// cont.startCont(types.ContainerStartOptions{})
	// listCont()
	// out, err := cli.ContainerLogs(context.Background(), cont.contID, types.ContainerLogsOptions{ShowStdout: true})
	// if err != nil {
	// 	return err
	// }

	// io.Copy(os.Stdout, out)
	// cont.stopCont(container.StopOptions{})
	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}
	fmt.Print(cli)
	rpc.Register(cont)
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(l, nil)
	// for {
	// 	// handle the connection

	// }
}
