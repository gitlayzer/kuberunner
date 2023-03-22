package type_secret

import corev1 "k8s.io/api/core/v1"

type SecretResp struct {
	Items []corev1.Secret `json:"items"`
	Total int             `json:"total"`
}

type SecretCreate struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Label     map[string]string `json:"label"`
	Data      map[string][]byte `json:"data"`
	Type      corev1.SecretType `json:"type"`
	Cluster   string            `json:"cluster"`
}
