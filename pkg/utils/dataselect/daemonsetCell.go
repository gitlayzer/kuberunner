package dataselect

import (
	appsv1 "k8s.io/api/apps/v1"
	"time"
)

type DaemonSetCell appsv1.DaemonSet

func (d DaemonSetCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}

func (d DaemonSetCell) GetName() string {
	return d.Name
}
