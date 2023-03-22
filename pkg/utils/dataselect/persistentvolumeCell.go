package dataselect

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type PersistentVolumeCell corev1.PersistentVolume

func (p PersistentVolumeCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p PersistentVolumeCell) GetName() string {
	return p.Name
}
