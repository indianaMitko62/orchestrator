package cluster

import (
	"fmt"
	"net/rpc"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

type ContainerNetworkConfig struct {
	NetworkID   string `yaml:"network_id,omitempty"`
	IPv4Address string `yaml:"ipv4_address,omitempty"`
	IPv6Address string `yaml:"ipv6_address,omitempty"`
	MACAddress  string `yaml:"mac_address,omitempty"`
	Gateway     string `yaml:"gateway,omitempty"`
	IPPrefixLen int    `yaml:"ip_prefix_len,omitempty"`
}

type ContainerConfig struct {
	Name            string                             `yaml:"name"`
	Image           string                             `yaml:"image"`
	Status          string                             `yaml:"status"`
	Volume          string                             `yaml:"volume,omitempty"`
	NetworkConfig   map[string]*ContainerNetworkConfig `yaml:"network_config"`
	Cmd             strslice.StrSlice                  `yaml:"cmd,omitempty"`
	Domainname      string                             `yaml:"domainname,omitempty"`
	WorkingDir      string                             `yaml:"working_dir,omitempty"`
	ExposedPorts    nat.PortSet                        `yaml:"exposed_ports,omitempty"`
	Healthcheck     container.HealthConfig             `yaml:"healthcheck,omitempty"`
	NetworkDisabled bool                               `yaml:"network_disabled,omitempty"`
	MacAddress      string                             `yaml:"mac_address,omitempty"`
	Privileged      bool                               `yaml:"privileged,omitempty"`
	ReadOnlyFS      bool                               `yaml:"read_only_fs,omitempty"`
	DNS             []string                           `yaml:"dns,omitempty"`

	PortBindings nat.PortMap `yaml:"port_bindings,omitempty"`
}

type NetworkConfig struct {
	NetworkID string `yaml:"network_id,omitempty"`
	Name      string `yaml:"name,omitempty"`
	Driver    string `yaml:"driver,omitempty"`
}

type ImageConfig struct {
	ImageID string `yaml:"id,omitempty"`
	Tag     string `yaml:"tag,omitempty"`
	Name    string `yaml:"name,omitempty"`
}

type VolumeConfig struct {
	VolumeID string `yaml:"volume_id,omitempty"`
	Driver   string `yaml:"driver,omitempty"`
}

type NodeSettings struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

type NodeManager struct {
	NodeSettings
	Client     *rpc.Client
	Containers map[string]*ContainerConfig `yaml:"containers"`
	Networks   map[string]*NetworkConfig   `yaml:"networks"`
	Volumes    map[string]*VolumeConfig    `yaml:"volumes"`
	Images     map[string]*ImageConfig     `yaml:"images"`
}

type ClusterState struct {
	Nodes map[string]*NodeManager `yaml:"nodes"`
}

func NewClusterState() *ClusterState {
	return &ClusterState{
		Nodes: make(map[string]*NodeManager),
	}
}

func (cs *ClusterState) CollectImages() { // probably won't be used in final version. Created for setup for node logic testing
	for _, node := range cs.Nodes {
		for _, cont := range node.Containers {
			parts := strings.Split(cont.Image, ":")
			var tag string
			name := parts[0]
			if len(parts) > 1 {
				tag = parts[1]
			}
			node.Images[cont.Image] = &ImageConfig{
				Name:    name,
				Tag:     tag,
				ImageID: "to be added via image inspection or by image creation in struct",
			}
		}
	}
}

func (n *NodeManager) Connect() error {
	client, err := rpc.DialHTTP("tcp", n.Address)
	if err != nil {
		return fmt.Errorf("could not connect to node's %s RPC service at %s: %w", n.Name, n.Address, err)
	}
	n.Client = client
	return nil
}
