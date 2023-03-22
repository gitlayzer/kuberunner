package dataselect

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type NodeCell corev1.Node

func (n NodeCell) GetCreation() time.Time {
	return n.CreationTimestamp.Time
}

func (n NodeCell) GetName() string {
	return n.Name
}
