package dataselect

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

type EventCell corev1.Event

func (e EventCell) GetCreation() time.Time {
	return e.CreationTimestamp.Time
}

func (e EventCell) GetName() string {
	return e.Name
}
