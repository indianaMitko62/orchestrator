package node

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/indianaMitko62/orchestrator/src/cluster"
	"gopkg.in/yaml.v3"
)

func (nsvc *NodeService) SendNodeStatus(URL string, nodeStatus *cluster.NodeStatus) error {
	NSToSend, _ := cluster.ToYaml(nodeStatus)
	fmt.Println("YAML Output:")
	yamlBytes := []byte(NSToSend)
	// fmt.Println(string(NSToSend))
	// fmt.Println(yamlBytes) // for testing
	req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(yamlBytes))
	if err != nil {
		nsvc.nodeLog.Logger.Error("Could not create POST request", "URL", URL)
		return err
	}
	req.Header.Set("nodeName", nsvc.Name)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		nsvc.nodeLog.Logger.Error("Could not send POST request to " + URL)
		return err
	}

	if resp.StatusCode == http.StatusOK {
		nsvc.nodeLog.Logger.Info("Node Status send successfully")
	}
	return nil
}

func (nsvc *NodeService) sendLogs(URL string, Log *cluster.Log) error {
	req, err := http.NewRequest(http.MethodPost, URL, Log.LogReader)
	if err != nil {
		nsvc.nodeLog.Logger.Error("Could not create POST request", "URL", URL)
		return err
	}
	req.Header.Set("nodeName", nsvc.Name)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		nsvc.nodeLog.Logger.Error("Could not send POST request to " + URL)
		return err
	}

	if resp.StatusCode == http.StatusOK {
		nsvc.nodeLog.Logger.Info("Cluster Change Outcome logs send successfully")
	}
	file, _ := os.Open(Log.FileName)
	file.Seek(-1, io.SeekEnd)
	Log.LogReader = file
	return nil
}

func (nsvc *NodeService) getClusterState(URL string) error {
	var cs cluster.ClusterState
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		nsvc.nodeLog.Logger.Error("Could not create POST request", "URL", URL)
		return err
	}
	req.Header.Set("nodeName", nsvc.Name)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("Could not send GET request to "+URL, "error", err)
		return err
	}
	if resp.StatusCode == http.StatusOK {
		yamlData, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading YAML data:", err)
			return err
		}
		fmt.Println(string(yamlData)) // for testing

		err = yaml.Unmarshal(yamlData, &cs)
		if err != nil {
			slog.Error("could not unmarshal cluster state yaml", "error", err)
			return err
		}
	} else {
		slog.Error("could not get cluster state", "URL", URL, "status", resp.Status)
	}
	*nsvc.DesiredNodeState = cs.Nodes[nsvc.Name].NodeState
	return nil
}
