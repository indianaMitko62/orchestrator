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
	NetworkID   string   `yaml:"network_id,omitempty"`
	IPv4Address string   `yaml:"ipv4_address,omitempty"`
	IPv6Address string   `yaml:"ipv6_address,omitempty"`
	MACAddress  string   `yaml:"mac_address,omitempty"`
	Gateway     string   `yaml:"gateway,omitempty"`
	DNS         []string `yaml:"dns,omitempty"`

	PortBindings nat.PortMap `yaml:"port_bindings,omitempty"`
}

type ContainerConfig struct {
	Name            string                 `yaml:"name"`
	Image           string                 `yaml:"image"`
	Status          string                 `yaml:"status"`
	Volume          string                 `yaml:"volume,omitempty"`
	NetworkConfig   ContainerNetworkConfig `yaml:"network_config,omitempty"`
	Cmd             strslice.StrSlice      `yaml:"cmd,omitempty"`
	Domainname      string                 `yaml:"domainname,omitempty"`
	WorkingDir      string                 `yaml:"working_dir,omitempty"`
	ExposedPorts    nat.PortSet            `yaml:"exposed_ports,omitempty"`
	Healthcheck     container.HealthConfig `yaml:"healthcheck,omitempty"`
	NetworkDisabled bool                   `yaml:"network_disabled,omitempty"`
	MacAddress      string                 `yaml:"mac_address,omitempty"`
	Privileged      bool                   `yaml:"privileged,omitempty"`
	ReadOnlyFS      bool                   `yaml:"read_only_fs,omitempty"`
}

type NetworkConfig struct {
	NetworkID string `yaml:"network_id,omitempty"`
	Name      string `yaml:"name,omitempty"`
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
	Containers []*ContainerConfig `yaml:"containers"`
	Networks   []*NetworkConfig   `yaml:"networks"`
	Volumes    []*VolumeConfig    `yaml:"volumes"`
	Images     []*ImageConfig     `yaml:"images"`
}

type ClusterState struct {
	Nodes map[string]*NodeManager `yaml:"nodes"`
}

func NewClusterState() *ClusterState {
	return &ClusterState{
		Nodes: make(map[string]*NodeManager),
	}
}

func (cs *ClusterState) CollectImages() {
	fmt.Println("alo da")
	fmt.Println(cs.Nodes)
	for _, node := range cs.Nodes {
		fmt.Println("alo da")
		var uniqueImages []string
		for _, cont := range node.Containers {
			uniqueImages = append(uniqueImages, cont.Image)
			fmt.Println(node)
		}
		for _, imageName := range uniqueImages {
			parts := strings.Split(imageName, ":")
			var tag string
			if len(parts) > 1 {
				tag = parts[1]
			}
			name := parts[0]

			imageConfig := &ImageConfig{
				Tag:  tag,
				Name: name,
			}
			node.Images = append(node.Images, imageConfig)
		}
	}
}

func (DesiredCS *ClusterState) Compare(CurrentCS *ClusterState) bool {

	return true
}

func (n *NodeManager) Connect() error {
	client, err := rpc.DialHTTP("tcp", n.Address)
	if err != nil {
		return fmt.Errorf("could not connect to node's %s RPC service at %s: %w", n.Name, n.Address, err)
	}
	n.Client = client
	return nil
}
