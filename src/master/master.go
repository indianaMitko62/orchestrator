package master

import (
	"net/rpc"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
)

type ContainerNetworkConfig struct {
	NetworkID   *string
	IPv4Address *string
	IPv6Address *string
	MACAddress  *string
	Gateway     *string

	DNS          []*string
	PortBindings nat.PortMap
}

type ContainerConfig struct {
	Name            *string
	Image           *string
	Volume          *string
	NetworkConfig   *ContainerNetworkConfig
	Cmd             *strslice.StrSlice
	Domainname      *string
	WorkingDir      *string
	ExposedPorts    *nat.PortSet
	Healthcheck     *container.HealthConfig
	NetworkDisabled *bool
	MacAddress      *string
	Privileged      *bool
	ReadOnlyFS      *bool
}

type NetworkConfig struct {
}

type VolumeConfig struct {
}

type NodeSettings struct {
	Name         *string
	Address      *string
	Containers   []*ContainerConfig
	Networks     []*NetworkConfig
	VolumeConfig []*VolumeConfig
}

type NodeManager struct {
	NodeSettings
	client *rpc.Client
}

type ClusterState struct {
	Nodes map[string]*NodeManager
}
