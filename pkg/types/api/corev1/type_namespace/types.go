package type_namespace

import corev1 "k8s.io/api/core/v1"

type NamespaceResp struct {
	Items []corev1.Namespace `json:"items"`
	Total int                `json:"total"`
}

type NamespaceCreate struct {
	Name    string            `json:"name"`
	Label   map[string]string `json:"label"`
	Cluster string            `json:"cluster"`
}
