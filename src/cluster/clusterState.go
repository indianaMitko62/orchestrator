package cluster

import (
	"strings"
)

type ClusterState struct {
	Nodes map[string]NodeState `yaml:"nodes"`
}

func NewClusterState() *ClusterState {
	return &ClusterState{
		Nodes: make(map[string]NodeState),
	}
}

func (cs *ClusterState) CollectImages() { // probably won't be used in final version. Created for setup for node logic testing
	for _, node := range cs.Nodes {
		for _, cont := range node.Containers {
			parts := strings.Split(cont.ContainerConfig.Image, ":")
			var tag string
			name := parts[0]
			if len(parts) > 1 {
				tag = parts[1]
			}
			//cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			node.Images[cont.ContainerConfig.Image] = &OrchImage{ //////////////////////////// to be checked again
				Name:          name,
				Tag:           tag,
				DesiredStatus: "pulled",
			}
		}
	}
}
