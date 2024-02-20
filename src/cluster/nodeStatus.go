package cluster

import (
	"io"
	"log/slog"
	"os"
)

type Log struct {
	Logger    *slog.Logger
	LogReader io.Reader
	FileName  string
}

func NewLog(file string) *Log {
	logFile, _ := os.Create(file) //not create, but open with append
	logReader, err := os.Open(file)
	if err != nil {
		slog.Error("Could not create logger")
	}
	logWriter := io.MultiWriter(os.Stdout, logFile)
	logger := slog.New(slog.NewTextHandler(logWriter, nil))
	return &Log{
		Logger:    logger,
		LogReader: logReader,
		FileName:  file,
	}
}

type NodeStatus struct {
	CPU              float32
	Memory           float32
	Disc             float32
	CurrentNodeState NodeState
}

func PostNodeStatus(URL string) (*NodeStatus, error) {
	var ns NodeStatus
	// resp, err := http.Get(URL)
	// if err != nil {
	// 	slog.Error("Could not send cluster state request to master", "error", err)
	// 	return nil, err
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode == http.StatusOK {
	// 	yamlData, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		slog.Error("Error reading YAML data:", err)
	// 		return &ClusterState{}, err
	// 	}
	// 	//fmt.Println(string(yamlData)) // for testing

	// 	err = yaml.Unmarshal(yamlData, &cs)
	// 	if err != nil {
	// 		slog.Error("could not unmarshal cluster state yaml", "error", err)
	// 		return &ClusterState{}, err
	// 	}
	// } else {
	// 	slog.Error("could not get cluster state", "URL", URL, "status", resp.Status)
	// }
	return &ns, nil
}

func (NS *NodeStatus) ToYaml() ([]byte, error) {
	// copyCS := *NS
	// yamlData, err := yaml.Marshal(copyCS)
	// if err != nil {
	// 	slog.Error("could create yaml representation")
	// }
	return []byte{}, nil
}
