package master

import (
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

type MasterService struct {
	CS                   *cluster.ClusterState
	ClusterChangeOutcome *cluster.ClusterChangeOutcome
}

func NewMasterService(cs *cluster.ClusterState) *MasterService {
	m := &MasterService{}
	m.CS = cs
	return m
}
