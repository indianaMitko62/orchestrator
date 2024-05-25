package master

import (
	"fmt"
	"io"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/indianaMitko62/orchestrator/src/cluster"
	"gopkg.in/yaml.v3"
)

func (msvc *MasterService) getClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	nodeName := r.Header.Get("nodeName")
	msvc.masterLog.Logger.Info("Received GET on /clusterState from", "node", nodeName, "IP", r.RemoteAddr)

	if !msvc.NodesStatus[nodeName].Operating {
		msvc.masterLog.Logger.Info("Node  online", "name", nodeName)
	}
	if msvc.CS == nil {
		w.Write([]byte("No cluster state configuration"))
		return
	}
	CSToSend, _ := cluster.ToYaml(msvc.CS)
	// fmt.Println("YAML Output:") // for testing
	// fmt.Println(string(CSToSend))

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Write(CSToSend)
}

func (msvc *MasterService) check_CS(cs *cluster.ClusterState) {
	// add default networking if configured
	if msvc.DefaultNetworking {
		// create macvlan network
		macvlanNetworkConfig := types.NetworkCreate{
			CheckDuplicate: true,
			Driver:         "macvlan",
			EnableIPv6:     false,
			IPAM: &network.IPAM{
				Driver: "default",
				Config: []network.IPAMConfig{
					{
						Subnet:  "192.168.42.0/24",
						Gateway: "192.168.42.1",
					},
				},
			},
			Options: map[string]string{
				"parent": "eth1",
			},
		}

		// connect containers to newly created default network
		node_number := 1
		for _, node := range cs.Nodes {
			node.Networks["default_container_network"] = &cluster.OrchNetwork{
				Name:          "default_container_network",
				DesiredStatus: "created",
				NetworkConfig: macvlanNetworkConfig,
			}
			cont_number := 0
			for _, cont := range node.Containers {
				host_address := fmt.Sprint(node_number) + fmt.Sprint(cont_number)
				if cont.NetworkingConfig == nil {
					cont.NetworkingConfig = &network.NetworkingConfig{}
				}
				if cont.NetworkingConfig.EndpointsConfig == nil {
					cont.NetworkingConfig.EndpointsConfig = map[string]*network.EndpointSettings{}
				}
				if cont.NetworkingConfig.EndpointsConfig["default_container_network"] == nil {
					cont.NetworkingConfig.EndpointsConfig["default_container_network"] = &network.EndpointSettings{
						NetworkID: "default_container_network",
						IPAMConfig: &network.EndpointIPAMConfig{
							IPv4Address: "192.168.42." + host_address,
						},
					}
				}
				cont_number++
			}
			node_number++
		}
	}
	// add default healthcheck if not not specified
}

func (msvc *MasterService) postClusterStateHandler(w http.ResponseWriter, r *http.Request) {
	nodeName := r.Header.Get("nodeName")
	msvc.masterLog.Logger.Info("Received POST on /clusterState from", "node", nodeName, "IP", r.RemoteAddr)

	var clusterState cluster.ClusterState
	yamlData, err := io.ReadAll(r.Body)
	if err != nil {
		msvc.masterLog.Logger.Error("Error reading YAML data:", err)
	}
	// fmt.Println(string(yamlData)) // for testing
	yaml.Unmarshal(yamlData, &clusterState)

	msvc.check_CS(&clusterState)

	msvc.CS = &clusterState
}
