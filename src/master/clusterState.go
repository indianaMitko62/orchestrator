package master

import (
	"net/rpc"

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

type VolumeConfig struct {
	VolumeID string `yaml:"volume_id,omitempty"`
	Driver   string `yaml:"driver,omitempty"`
}

type NodeSettings struct {
	Name         string             `yaml:"name"`
	Address      string             `yaml:"address"`
	Containers   []*ContainerConfig `yaml:"containers"`
	Networks     []*NetworkConfig   `yaml:"networks"`
	VolumeConfig []*VolumeConfig    `yaml:"volume_config"`
}

type NodeManager struct {
	NodeSettings
	Client *rpc.Client
}

type ClusterState struct {
	Nodes map[string]*NodeManager `yaml:"nodes"`
}

func NewClusterState() *ClusterState {
	return &ClusterState{
		Nodes: make(map[string]*NodeManager),
	}
}
