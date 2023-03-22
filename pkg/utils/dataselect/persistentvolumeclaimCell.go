package dataselect

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type PersistentVolumeClaimCell corev1.PersistentVolumeClaim

func (p PersistentVolumeClaimCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p PersistentVolumeClaimCell) GetName() string {
	return p.Name
}
