package node

import (
	"github.com/docker/docker/api/types"
	"github.com/indianaMitko62/orchestrator/src/cluster"
)

func (nsvc *NodeService) DeployNewImage(img *cluster.OrchImage) {
	img.Cli = nsvc.cli
	_, err := img.PullImg(&types.ImagePullOptions{
		All:           img.All,
		RegistryAuth:  img.RegistryAuth,
		Platform:      img.Platform,
		PrivilegeFunc: nil,
	})
	if img.CurrentStatus == img.DesiredStatus && err == nil {
		nsvc.CurrentNodeState.Images[img.Name] = img
		nsvc.clusterChangeLog.Logger.Info("Image successfully pulled", "name", img.Name)
	} else {
		nsvc.clusterChangeLog.Logger.Info("Could not pull image", "name", img.Name, "err", err.Error())
	}
}
