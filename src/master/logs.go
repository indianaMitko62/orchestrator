package master

import (
	"io"
	"net/http"
	"os"

	"github.com/indianaMitko62/orchestrator/src/cluster"
)

func (msvc *MasterService) postLogsHandler(w http.ResponseWriter, r *http.Request) {
	nodeName := r.Header.Get("nodeName")
	msvc.masterLog.Logger.Info("recieved POST on /logs from", "name", nodeName, "IP", r.RemoteAddr)
	_, ok := msvc.CS.Nodes[nodeName]
	if !ok {
		newNode := *cluster.NewNodeState()
		newNode.Active = false
		if msvc.CS.Nodes == nil {
			msvc.masterLog.Logger.Error("Cluster State not present. Check YAML configuration")
		}
		msvc.CS.Nodes[nodeName] = newNode
		msvc.masterLog.Logger.Info("added Node to CS", "name", nodeName, "IP", r.RemoteAddr)
	}
	logData, err := io.ReadAll(r.Body)
	if err != nil {
		msvc.masterLog.Logger.Error("Error reading YAML data:", err)
	}
	// fmt.Println(string(logData)) // for testing
	f, err := os.OpenFile(msvc.LogsPath+nodeName+"Logs", os.O_APPEND|os.O_WRONLY, 0600) // separate go routine
	if err != nil {
		msvc.masterLog.Logger.Warn("Could not open file. Trying to create it", "name", msvc.LogsPath+nodeName+"Logs")
		_, err := os.Create(msvc.LogsPath + nodeName + "Logs")
		if err != nil {
			msvc.masterLog.Logger.Error("Could not create file", "name", msvc.LogsPath+nodeName+"Logs", "error", err)
			return
		}
		f, err = os.OpenFile(msvc.LogsPath+nodeName+"Logs", os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			msvc.masterLog.Logger.Error("Could not open file", "name", msvc.LogsPath+nodeName+"Logs", "error", err)
			return
		}
	}
	f.Write(logData)
	f.Close()
	msvc.masterLog.Logger.Info("Logs written", "file", msvc.LogsPath+nodeName+"Logs")
}
