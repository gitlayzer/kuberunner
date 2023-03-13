package dataselect

import (
	v1 "k8s.io/api/storage/v1"
	"time"
)

type StorageClassCell v1.StorageClass

func (s StorageClassCell) GetCreation() time.Time {
	return s.CreationTimestamp.Time
}

func (s StorageClassCell) GetName() string {
	return s.Name
}
