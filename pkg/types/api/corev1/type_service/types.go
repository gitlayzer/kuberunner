package type_service

import corev1 "k8s.io/api/core/v1"

type ServiceResp struct {
	Items []corev1.Service `json:"items"`
	Total int              `json:"total"`
}

type ServiceCreate struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Type       string            `json:"type"`
	Label      map[string]string `json:"label"`
	Selector   map[string]string `json:"selector"`
	PortsName  string            `json:"ports_name"`
	PortsPort  int32             `json:"ports_port"`
	Protocol   string            `json:"protocol"`
	TargetPort int32             `json:"target_port"`
	Cluster    string            `json:"cluster"`
}
