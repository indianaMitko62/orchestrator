package master

import (
	"github.com/docker/docker/client"
)

type MasterService struct {
	cli *client.Client

	CS *ClusterState
}

func NewMasterService(cs *ClusterState) *MasterService {
	m := &MasterService{}
	m.CS = cs
	return m
}
