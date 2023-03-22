package type_pod

import corev1 "k8s.io/api/core/v1"

type PodsResp struct {
	Items []corev1.Pod `json:"items"`
	Total int          `json:"total"`
}
