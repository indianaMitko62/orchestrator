package master

import (
	"net/http"

	"github.com/indianaMitko62/orchestrator/src/cluster"
)

func (msvc *MasterService) getClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	nodeName := r.Header.Get("nodeName")
	msvc.masterLog.Logger.Info("Received GET on /clusterState from", "node", nodeName, "IP", r.RemoteAddr)

	if !msvc.NodesStatus[nodeName].Active {
		msvc.masterLog.Logger.Info("Inactive node back online", "name", nodeName)
	}
	CSToSend, _ := cluster.ToYaml(msvc.CS)
	// fmt.Println("YAML Output:") // for testing
	// fmt.Println(string(CSToSend))

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(CSToSend)
}
