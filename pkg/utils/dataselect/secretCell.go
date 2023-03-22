package dataselect

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type SecretCell corev1.Secret

func (s SecretCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s SecretCell) GetName() string {
	return s.Name
}
