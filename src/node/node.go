package node

import (
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"

	"github.com/indianaMitko62/orchestrator/src/cluster"
	"gopkg.in/yaml.v3"
)

/*
TODO: functions managing overall node performance and loading(cpu, memory, disk)
*/
var desiredClusterState cluster.ClusterState
var currentClusterState cluster.ClusterState

func (ns *NodeService) Node() error {
	masterURL := "http://localhost:1986/clusterState"
	for i := 0; i < 2; i++ {
		fmt.Println(i)
		resp, err := http.Get(masterURL)
		if err != nil {
			fmt.Println("Could not send cluster state request to master", "error", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			yamlData, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("Error reading YAML data:", err)
				return err
			}
			fmt.Println(string(yamlData))

			err = yaml.Unmarshal(yamlData, &desiredClusterState)
			if err != nil {
				slog.Error("could not unmarshall cluster state yaml", "error", err)
				return err
			}
		} else {
			slog.Error("could not get cluster state.", "status", resp.Status)
		}

		if maps.Equal(desiredClusterState.Nodes, currentClusterState.Nodes) {

		}

		currentClusterState = desiredClusterState
	}
	fmt.Println("out")
	return nil
}
