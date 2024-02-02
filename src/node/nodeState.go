package node

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/indianaMitko62/orchestrator/src/cluster"
	"gopkg.in/yaml.v3"
)

type NodeState struct {
	Containers map[string]*OrchContainer
	Networks   map[string]*OrchNetwork
	Images     map[string]*OrchImage
	Volumes    map[string]*OrchVolume
}

func NewNodeState() *NodeState {
	return &NodeState{
		Containers: make(map[string]*OrchContainer),
		Networks:   make(map[string]*OrchNetwork),
		Images:     make(map[string]*OrchImage),
		Volumes:    make(map[string]*OrchVolume),
	}
}

/*
TODO: functions managing overall node performance and loading(cpu, memory, disk) and overall node logic
*/
var recievedClusterState *cluster.ClusterState
var desiredClusterState *cluster.ClusterState // to be removed. Stays for now to test clusterState receive
// Or not to be removed. Can be used for quick comparison to avoid useless node state parsing

func getClusterState(URL string) (*cluster.ClusterState, error) {
	var cs cluster.ClusterState
	resp, err := http.Get(URL)
	if err != nil {
		slog.Error("Could not send cluster state request to master", "error", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		yamlData, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading YAML data:", err)
			return &cluster.ClusterState{}, err
		}
		fmt.Println(string(yamlData)) // for testing

		err = yaml.Unmarshal(yamlData, &cs)
		if err != nil {
			slog.Error("could not unmarshal cluster state yaml", "error", err)
			return &cluster.ClusterState{}, err
		}
	} else {
		slog.Error("could not get cluster state", "URL", URL, "status", resp.Status)
	}
	return &cs, nil
}

func (nsvc *NodeService) parseContainers(containers map[string]*cluster.ContainerConfig) error {
	for _, cont := range containers {
		parsedContainer := &OrchContainer{
			ContStatus: cont.Status,
			ContainerConfig: &container.Config{
				Hostname:        cont.Name,
				Image:           cont.Image,
				Domainname:      cont.Domainname,
				ExposedPorts:    cont.ExposedPorts,
				Cmd:             cont.Cmd,
				Healthcheck:     &cont.Healthcheck,
				WorkingDir:      cont.WorkingDir,
				NetworkDisabled: cont.NetworkDisabled,
				MacAddress:      cont.MacAddress,
			},
			HostConfig: &container.HostConfig{
				ReadonlyRootfs: cont.ReadOnlyFS,
				PortBindings:   cont.PortBindings,
				DNS:            cont.DNS,
				Privileged:     cont.Privileged,
			},
			NetworkingConfig: &network.NetworkingConfig{
				EndpointsConfig: map[string]*network.EndpointSettings{},
			},
		}
		for name, setting := range cont.NetworkConfig {
			parsedContainer.NetworkingConfig.EndpointsConfig[name] = &network.EndpointSettings{
				NetworkID:   setting.NetworkID,
				IPAddress:   setting.IPv4Address,
				IPPrefixLen: setting.IPPrefixLen,
				MacAddress:  setting.MACAddress,
			}
		}
		nsvc.DesiredNodeState.Containers[parsedContainer.ContainerConfig.Hostname] = parsedContainer
		nsvc.CreateCont(parsedContainer)
	}
	return nil
}

func (nsvc *NodeService) parseImages(images map[string]*cluster.ImageConfig) error {
	for name, img := range images {
		fmt.Println(name, img)
	}
	return nil
}

func (nsvc *NodeService) parseNetworks(networks map[string]*cluster.NetworkConfig) error {
	for name, net := range networks {
		fmt.Println(name, net)
	}
	return nil
}

func (nsvc *NodeService) parseVolumes(volumes map[string]*cluster.VolumeConfig) error {
	for name, vol := range volumes {
		fmt.Println(name, vol)
	}
	return nil
}

func (nsvc *NodeService) getNodeState(CS *cluster.ClusterState) error {
	fmt.Println("aaaa" + nsvc.Name)
	thisNodeManager := CS.Nodes[nsvc.Name]
	fmt.Println(thisNodeManager)

	nsvc.parseContainers(thisNodeManager.Containers)
	nsvc.parseImages(thisNodeManager.Images)
	nsvc.parseNetworks(thisNodeManager.Networks)
	nsvc.parseVolumes(thisNodeManager.Volumes)
	return nil
}
