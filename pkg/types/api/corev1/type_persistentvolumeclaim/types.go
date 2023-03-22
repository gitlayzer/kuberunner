package type_persistentvolumeclaim

import corev1 "k8s.io/api/core/v1"

type PersistentVolumeClaimResp struct {
	Items []corev1.PersistentVolumeClaim `json:"items"`
	Total int                            `json:"total"`
}

type PersistentVolumeClaimCreate struct {
	Name                    string `json:"name"`
	Namespace               string `json:"namespace"`
	Label                   map[string]string
	AccessMode              string            `json:"access_mode"`
	Storage                 string            `json:"storage"`
	StorageClass            string            `json:"storage_class"`
	BindingPersistentVolume map[string]string `json:"binding_persistent_volume"`
	Cluster                 string            `json:"cluster"`
}
