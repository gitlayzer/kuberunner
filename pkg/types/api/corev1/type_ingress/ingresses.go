package type_ingress

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/utils/dataselect"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Ingress ingress

type ingress struct{}

func (i *ingress) GetIngresses(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (ingressResp *IngressResp, err error) {
	ingressList, err := client.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Ingress列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: i.toCells(ingressList.Items),
		DataSelectorQuery: &dataselect.DataSelectorQuery{
			FilterQuery: &dataselect.FilterQuery{
				Name: filterName,
			},
			Pagination: &dataselect.PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	filtered := selectableData.Filter()

	total := len(filtered.GenericDataList)

	data := filtered.Sort().Paginate()

	ingresses := i.fromCells(data.GenericDataList)

	return &IngressResp{
		Items: ingresses,
		Total: total,
	}, nil
}

func (i *ingress) GetIngressDetail(client *kubernetes.Clientset, ingressName, namespace string) (service *v1.Ingress, err error) {
	ingresses, err := client.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Ingress详情失败: %v", err))
	}

	return ingresses, nil
}

func (i *ingress) UpdateIngress(client *kubernetes.Clientset, namespace, content string) (err error) {

	var ingress v1.Ingress

	err = json.Unmarshal([]byte(content), &ingress)
	if err != nil {
		return errors.New(fmt.Sprintf("反序列化失败: %v\n", err))
	}

	_, err = client.NetworkingV1().Ingresses(namespace).Update(context.TODO(), &ingress, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("更新Ingress失败: %v", err))
	}

	return nil
}

func (i *ingress) CreateIngress(client *kubernetes.Clientset, ingressCreate *IngressCreate) (err error) {
	ingressRules := make([]v1.IngressRule, 0, len(ingressCreate.Hosts))
	ingressPaths := make([]v1.HTTPIngressPath, 0)

	for _, ingressHTTPPaths := range ingressCreate.Hosts {
		host := ingressHTTPPaths.Host
		for _, ingressHttpPath := range ingressHTTPPaths.HttpPaths {
			ingressRules = append(ingressRules, v1.IngressRule{
				Host: host,
				IngressRuleValue: v1.IngressRuleValue{
					HTTP: &v1.HTTPIngressRuleValue{
						Paths: []v1.HTTPIngressPath{
							{
								Path:     ingressHttpPath.Path,
								PathType: &ingressHttpPath.PathType,
								Backend: v1.IngressBackend{
									Service: &v1.IngressServiceBackend{
										Name: ingressHttpPath.ServiceName,
										Port: v1.ServiceBackendPort{
											Number: ingressHttpPath.ServicePort,
										},
									},
								},
							},
						},
					},
				},
			})
			ingressPaths = append(ingressPaths, v1.HTTPIngressPath{
				Path:     ingressHttpPath.Path,
				PathType: &ingressHttpPath.PathType,
				Backend: v1.IngressBackend{
					Service: &v1.IngressServiceBackend{
						Name: ingressHttpPath.ServiceName,
						Port: v1.ServiceBackendPort{
							Number: ingressHttpPath.ServicePort,
						},
					},
				},
			})
		}
	}

	ingresses := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ingressCreate.Name,
			Namespace: ingressCreate.Namespace,
			Labels:    ingressCreate.Label,
		},
		Spec: v1.IngressSpec{
			IngressClassName: &ingressCreate.IngressClassName,
			Rules:            ingressRules,
		},
	}

	_, err = client.NetworkingV1().Ingresses(ingressCreate.Namespace).Create(context.TODO(), ingresses, metav1.CreateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("创建Ingress失败: %v", err))
	}

	return nil
}

func (i *ingress) DeleteIngress(client *kubernetes.Clientset, ingressName, namespace string) (err error) {
	if err := client.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{}); err != nil {
		return errors.New(fmt.Sprintf("删除Ingress失败: %v", err))
	}

	return nil
}

func (i *ingress) toCells(std []v1.Ingress) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for j := range std {
		cells[j] = dataselect.IngressCell(std[j])
	}
	return cells
}

func (i *ingress) fromCells(cells []dataselect.DataCell) []v1.Ingress {
	ingresses := make([]v1.Ingress, len(cells))
	for j := range cells {
		ingresses[j] = v1.Ingress(cells[j].(dataselect.IngressCell))
	}
	return ingresses
}
