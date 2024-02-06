package cluster

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/rpc"
	"strings"

	"gopkg.in/yaml.v3"
)

type NodeSettings struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

type NodeManager struct {
	NodeSettings
	Client *rpc.Client
	NodeState
}

type ClusterState struct {
	Nodes map[string]*NodeManager `yaml:"nodes"`
}

func NewClusterState() *ClusterState {
	return &ClusterState{
		Nodes: make(map[string]*NodeManager),
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
			node.Images[cont.ContainerConfig.Image] = &OrchImage{ //////////////////////////// to be checked again
				Name: name,
				Tag:  tag,
			}
		}
	}
}

func (n *NodeManager) Connect() error {
	client, err := rpc.DialHTTP("tcp", n.Address)
	if err != nil {
		return fmt.Errorf("could not connect to node's %s RPC service at %s: %w", n.Name, n.Address, err)
	}
	n.Client = client
	return nil
}

func GetClusterState(URL string) (*ClusterState, error) {
	var cs ClusterState
	resp, err := http.Get(URL)
	if err != nil {
		slog.Error("Could not send cluster state request to master", "error", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		yamlData, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading YAML data:", err)
			return &ClusterState{}, err
		}
		fmt.Println(string(yamlData)) // for testing

		err = yaml.Unmarshal(yamlData, &cs)
		if err != nil {
			slog.Error("could not unmarshal cluster state yaml", "error", err)
			return &ClusterState{}, err
		}
	} else {
		slog.Error("could not get cluster state", "URL", URL, "status", resp.Status)
	}
	return &cs, nil
}
