package type_configmap

import corev1 "k8s.io/api/core/v1"

type ConfigMapResp struct {
	Items []corev1.ConfigMap `json:"items"`
	Total int                `json:"total"`
}

type ConfigMapCreate struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Data      map[string]string `json:"data"`
	Cluster   string            `json:"cluster"`
}
