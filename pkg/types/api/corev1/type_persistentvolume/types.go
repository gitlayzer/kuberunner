package type_persistentvolume

import corev1 "k8s.io/api/core/v1"

type PersistentVolumeResp struct {
	Items []corev1.PersistentVolume `json:"items"`
	Total int                       `json:"total"`
}

type PersistentVolumeCreate struct {
	Name                          string            `json:"name"`
	Label                         map[string]string `json:"label"`
	StorageClass                  string            `json:"storage_class"`
	Storage                       string            `json:"storage"`
	AccessMode                    string            `json:"access_mode"`
	PersistentVolumeReclaimPolicy string            `json:"persistent_volume_reclaim_policy"`
	Nfs                           NfsCreate         `json:"nfs,omitempty"`
	Rbd                           CephRbdCreate     `json:"rbd,omitempty"`
	HostPath                      HostPathCreate    `json:"hostPath,omitempty"`
	Local                         LocalCreate       `json:"local,omitempty"`
	CephFs                        CephFsCreate      `json:"cephfs,omitempty"`
	Cluster                       string            `json:"cluster"`
}

type LocalCreate struct {
	Path   string  `json:"path"`
	FsType *string `json:"fsType"`
}

type NfsCreate struct {
	Path   string `json:"path"`
	Server string `json:"server"`
}

type CephRbdCreate struct {
	Monitors  []string               `json:"monitors"`
	Pool      string                 `json:"pool"`
	Image     string                 `json:"image"`
	User      string                 `json:"user"`
	SecretRef corev1.SecretReference `json:"secret_ref"`
	ReadOnly  bool                   `json:"read_only"`
	FsType    string                 `json:"fsType"`
}

type HostPathCreate struct {
	Path string `json:"path"`
}

type CephFsCreate struct {
	Monitors  []string               `json:"monitors"`
	Path      string                 `json:"path"`
	User      string                 `json:"user"`
	SecretRef corev1.SecretReference `json:"secret_ref"`
}
