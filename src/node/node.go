package node

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

/*
TODO: functions managing overall node performance and loading(cpu, memory, disk)
*/

func (ns *NodeService) Node() {
	masterURL := "http://localhost:1986/clusterState"
	slog.Info("started")
	resp, err := http.Get(masterURL)
	if err != nil {
		fmt.Println("Error sending request to master:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		yamlData, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading YAML data:", err)
			return
		}
		fmt.Println(string(yamlData))

		// var clusterState ClusterState
		// err = yaml.Unmarshal(yamlData, &clusterState)
		// if err != nil {
		// 	fmt.Println("Error unmarshaling YAML data:", err)
		// 	return
		// }

		// fmt.Println("Received YAML data:", clusterState)
	} else {
		fmt.Println("Error getting YAML data. Status:", resp.Status)
	}
}
