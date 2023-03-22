package type_ingress

import v1 "k8s.io/api/networking/v1"

type IngressResp struct {
	Items []v1.Ingress `json:"items"`
	Total int          `json:"total"`
}

type IngressCreate struct {
	Name             string              `json:"name"`
	Namespace        string              `json:"namespace"`
	Label            map[string]string   `json:"label"`
	Hosts            []*IngressHTTPPaths `json:"hosts"`
	Cluster          string              `json:"cluster"`
	IngressClassName string              `json:"ingress_class_name"`
}

type IngressHTTPPaths struct {
	Host      string             `json:"host"`
	HttpPaths []*IngressHttpPath `json:"http_paths"`
}

type IngressHttpPath struct {
	Path        string      `json:"path"`
	PathType    v1.PathType `json:"path_type"`
	ServiceName string      `json:"service_name"`
	ServicePort int32       `json:"service_port"`
}
