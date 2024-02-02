package node

import (
	"fmt"
	"maps"
)

func (nsvc *NodeService) Node() error {
	clusterStateURL := "http://localhost:1986/clusterState"
	for i := 0; i < 1; i++ {
		recievedClusterState, _ = getClusterState(clusterStateURL)
		if desiredClusterState == nil {
			nsvc.getNodeState(recievedClusterState)
			desiredClusterState = recievedClusterState
		} else if !maps.Equal(recievedClusterState.Nodes, desiredClusterState.Nodes) {
			fmt.Print(false)
			nsvc.getNodeState(recievedClusterState)
			desiredClusterState = recievedClusterState // probably to be moved; depends on exact desiredClusterState usage later on
		}
	}
	fmt.Println("out")
	return nil
}
