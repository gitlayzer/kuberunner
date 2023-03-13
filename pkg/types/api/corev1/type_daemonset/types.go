package type_daemonset

import appsv1 "k8s.io/api/apps/v1"

type DaemonSetResp struct {
	Items []appsv1.DaemonSet `json:"items"`
	Total int                `json:"total"`
}

type DaemonSetCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
	Cluster       string            `json:"cluster"`
}
