package dataselect

import (
	appsv1 "k8s.io/api/apps/v1"
	"time"
)

type DeploymentCell appsv1.Deployment

func (d DeploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d DeploymentCell) GetName() string {
	return d.Name
}
