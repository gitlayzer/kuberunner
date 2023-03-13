package type_storageclass

import v1 "k8s.io/api/storage/v1"

type StorageClassResp struct {
	Items []v1.StorageClass `json:"items"`
	Total int               `json:"total"`
}
