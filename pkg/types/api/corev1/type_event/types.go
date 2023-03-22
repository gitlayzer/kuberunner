package type_event

import corev1 "k8s.io/api/core/v1"

type EventResp struct {
	Items []corev1.Event `json:"items"`
	Total int            `json:"total"`
}
