package master

import (
	"io"
	"net/http"

	"github.com/indianaMitko62/orchestrator/src/cluster"
	"gopkg.in/yaml.v3"
)

func (msvc *MasterService) getClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	nodeName := r.Header.Get("nodeName")
	msvc.masterLog.Logger.Info("Received GET on /clusterState from", "node", nodeName, "IP", r.RemoteAddr)

	if !msvc.NodesStatus[nodeName].Active {
		msvc.masterLog.Logger.Info("Node  online", "name", nodeName)
	}
	if msvc.CS == nil {
		w.Write([]byte("No cluster state configuration"))
		return
	}
	CSToSend, _ := cluster.ToYaml(msvc.CS)
	// fmt.Println("YAML Output:") // for testing
	// fmt.Println(string(CSToSend))

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(CSToSend)
}

func (msvc *MasterService) postClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	nodeName := r.Header.Get("nodeName")
	msvc.masterLog.Logger.Info("Received POST on /clusterState from", "node", nodeName, "IP", r.RemoteAddr)

	var clusterState cluster.ClusterState
	yamlData, err := io.ReadAll(r.Body)
	if err != nil {
		msvc.masterLog.Logger.Error("Error reading YAML data:", err)
	}
	// fmt.Println(string(yamlData)) // for testing
	yaml.Unmarshal(yamlData, &clusterState)

	msvc.CS = &clusterState
}
