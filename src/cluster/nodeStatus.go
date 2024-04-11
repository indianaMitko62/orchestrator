package cluster

import (
	"io"
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v3"
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
		slog.Error("Could not create logger" + file)
		return nil
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
	CPU              float64
	Memory           float64
	Disk             float64
	CurrentNodeState NodeState
	Operating        bool
	Timestamp        time.Time
}

func ToYaml(data interface{}) ([]byte, error) {
	dataCopy := data
	yamlData, err := yaml.Marshal(dataCopy)
	if err != nil {
		slog.Error("could create yaml representation")
	}
	return yamlData, nil
}
