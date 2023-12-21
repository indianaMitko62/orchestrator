package main

import (
	"log"
	"net/rpc"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type Args struct {
	Cont Container
}

type Container struct {
	ContainerConfig  container.Config
	Image_name       string
	HostConfig       *container.HostConfig
	networkingConfig *network.NetworkingConfig
	ContainerName    string
	ContID           string
}

func main() {

	client1, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	cont := Container{
		container.Config{
			Image: "alpine:latest",
			Cmd:   []string{"echo", "maikatibe"},
		},
		"alpine", nil, nil, "cont1", "", ""}
	var reply string
	args := Args{Cont: cont}

	// err = client1.Call("Container.NewClient", args, &reply)
	// if err != nil {
	// 	log.Fatal("Client :", err)
	// }
	//args := &Args{types.ImagePullOptions{}}
	err = client1.Call("Container.CreateCont", args, &reply)
	if err != nil {
		log.Fatal("create error:", err)
	}
	//fmt.Printf("CreateCont", args.Opts)

}
