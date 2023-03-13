package type_node

import corev1 "k8s.io/api/core/v1"

type NodeResp struct {
	Items []corev1.Node `json:"items"`
	Total int           `json:"total"`
}
