package master

import (
	"fmt"
	"io"
	"net/http"
	"sort"

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
		ns := msvc.NodesStatus[nodeName]
		ns.Operating = true
		msvc.NodesStatus[nodeName] = ns
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

		cont_number := 0
		nodeNames := make([]string, 0)
		for k := range cs.Nodes {
			nodeNames = append(nodeNames, k)
		}
		sort.Strings(nodeNames)
		for _, nodeName := range nodeNames {
			node := cs.Nodes[nodeName]
			node.Networks["default_container_network"] = &cluster.OrchNetwork{
				Name:          "default_container_network",
				DesiredStatus: "created",
				NetworkConfig: macvlanNetworkConfig,
			}
			// contNames := make([]string, 0)
			// for k := range node.Containers {
			// 	contNames = append(contNames, k)
			// }
			// sort.Strings(contNames)
			// for _, name := range contNames {
			for name, cont := range node.Containers {
				// cont := node.Containers[name]
				host_address := fmt.Sprint(cont_number + 20)
				IP := "192.168.42." + host_address
				if msvc.nodeNameToIP[name] == "" {
					msvc.nodeNameToIP[name] = IP
					fmt.Println(name + " : " + IP)
				}
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
							IPv4Address: msvc.nodeNameToIP[name],
						},
					}
				}
				cont_number++
			}
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
	yaml.Unmarshal(yamlData, &clusterState)

	// fmt.Println(string(yamlData))

	fmt.Println(clusterState)
	for name := range clusterState.Nodes {
		fmt.Println(name)
	}
	msvc.check_CS(&clusterState)
	msvc.CS = &clusterState
}
