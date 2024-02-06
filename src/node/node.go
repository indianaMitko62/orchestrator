package node

import (
	"fmt"
	"log/slog"

	"github.com/indianaMitko62/orchestrator/src/cluster"
)

/*
TODO: functions managing overall node performance and loading(cpu, memory, disk) and overall node logic
*/

func (nsvc *NodeService) Node() error {
	clusterStateURL := "http://localhost:1986/clusterState"
	for i := 0; i < 1; i++ {
		recievedClusterState, _ := cluster.GetClusterState(clusterStateURL)
		fmt.Println(recievedClusterState)
		nsvc.DesiredNodeState = &recievedClusterState.Nodes[nsvc.Name].NodeState
		slog.Info("in getNodeState", "container name", nsvc.DesiredNodeState.Containers["Container1"].ContainerConfig.Hostname) // testing purposes
	}
	fmt.Println("out")
	return nil
}
