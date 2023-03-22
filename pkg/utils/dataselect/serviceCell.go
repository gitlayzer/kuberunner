package dataselect

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type ServiceCell corev1.Service

func (s ServiceCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s ServiceCell) GetName() string {
	return s.Name
}
