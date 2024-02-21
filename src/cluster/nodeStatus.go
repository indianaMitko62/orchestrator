package cluster

import (
	"io"
	"log/slog"
	"os"

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

func ToYaml(data interface{}) ([]byte, error) {
	dataCopy := data
	yamlData, err := yaml.Marshal(dataCopy)
	if err != nil {
		slog.Error("could create yaml representation")
	}
	return yamlData, nil
}
