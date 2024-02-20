package cluster

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"gopkg.in/yaml.v3"
)

type MasterServiceConfig struct {
}

type NodeSettings struct {
	Name             string `yaml:"name"`
	Address          string `yaml:"address"`
	MasterAddress    string `yaml:"master_address"`
	Port             string `yaml:"port"`
	LogsPath         string `yaml:"logspath"`
	ClusterStatePath string `yaml:"clusterchangepath"`
}

type NodeManager struct {
	NodeSettings
	NodeState
}

type ClusterState struct {
	Nodes map[string]*NodeManager `yaml:"nodes"`
}

type ClusterChangeOutcome struct {
	Successful bool
	Logs       map[string]string
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
			//cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			node.Images[cont.ContainerConfig.Image] = &OrchImage{ //////////////////////////// to be checked again
				Name:          name,
				Tag:           tag,
				DesiredStatus: "pulled",
			}
		}
	}
}

func GetClusterState(URL string) (*ClusterState, error) {
	var cs ClusterState
	resp, err := http.Get(URL)
	if err != nil {
		slog.Error("Could not send cluster state request to master", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		yamlData, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading YAML data:", err)
			return &ClusterState{}, err
		}
		//fmt.Println(string(yamlData)) // for testing

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

func (CS *ClusterState) ToYaml() ([]byte, error) {
	copyCS := *CS
	yamlData, err := yaml.Marshal(copyCS)
	if err != nil {
		slog.Error("could create yaml representation")
	}
	return yamlData, nil
}
