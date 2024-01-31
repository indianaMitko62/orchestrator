package master

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"gopkg.in/yaml.v3"
)

func (CS *ClusterState) ToYaml() ([]byte, error) {
	copyCS := *CS
	yamlData, err := yaml.Marshal(copyCS)
	if err != nil {
		slog.Error("could create yaml representation")
	}
	return yamlData, nil
}

func (CS *ClusterState) HandleClusterState(w http.ResponseWriter, r *http.Request) {
	nodeIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		slog.Error("Error parsing node IP:", err)
		return
	}

	var nodeName string // node is alive???
	for name, node := range CS.Nodes {
		if node.Address == nodeIP {
			nodeName = name
			break
		}
	}

	slog.Info("Received cluster state request", "node", nodeName, "IP", nodeIP)

	CSToSend, _ := CS.ToYaml()
	fmt.Println("YAML Output:")
	fmt.Println(string(CSToSend))

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(CSToSend)
}
