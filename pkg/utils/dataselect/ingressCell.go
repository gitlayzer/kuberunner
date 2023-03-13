package dataselect

import (
	v1 "k8s.io/api/networking/v1"
	"time"
)

type IngressCell v1.Ingress

func (i IngressCell) GetCreation() time.Time {
	return i.CreationTimestamp.Time
}

func (i IngressCell) GetName() string {
	return i.Name
}
