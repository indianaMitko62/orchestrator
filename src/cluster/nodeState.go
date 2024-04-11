package cluster

type NodeState struct {
	Containers map[string]*OrchContainer `yaml:"containers"`
	Networks   map[string]*OrchNetwork   `yaml:"networks"`
	Volumes    map[string]*OrchVolume    `yaml:"volumes"`
	Images     map[string]*OrchImage     `yaml:"images"`
	Active     bool                      `yaml:"active"`
}

func NewNodeState() *NodeState {
	return &NodeState{
		Containers: make(map[string]*OrchContainer),
		Networks:   make(map[string]*OrchNetwork),
		Images:     make(map[string]*OrchImage),
		Volumes:    make(map[string]*OrchVolume),
		Active:     true,
	}
}
