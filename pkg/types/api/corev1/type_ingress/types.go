package type_ingress

import v1 "k8s.io/api/networking/v1"

type IngressResp struct {
	Items []v1.Ingress `json:"items"`
	Total int          `json:"total"`
}

type IngressCreate struct {
	Name             string               `json:"name"`
	Namespace        string               `json:"namespace"`
	Label            map[string]string    `json:"label"`
	Hosts            map[string]*HttpPath `json:"hosts"`
	Cluster          string               `json:"cluster"`
	IngressClassName string               `json:"ingress_class_name"`
}

type HttpPath struct {
	Path        string      `json:"path"`
	PathType    v1.PathType `json:"path_type"`
	ServiceName string      `json:"service_name"`
	ServicePort int32       `json:"service_port"`
}
