package type_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/utils/dataselect"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

var Services services

type services struct{}

func (s *services) GetServices(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (serviceResp *ServiceResp, err error) {
	serviceList, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Service列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: s.toCells(serviceList.Items),
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

	services := s.fromCells(data.GenericDataList)

	return &ServiceResp{
		Items: services,
		Total: total,
	}, nil
}

func (s *services) GetServiceDetail(client *kubernetes.Clientset, serviceName, namespace string) (service *corev1.Service, err error) {
	service, err = client.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Service详情失败: %v", err))
	}

	return service, nil
}

func (s *services) UpdateService(client *kubernetes.Clientset, namespace, content string) (err error) {

	var service = &corev1.Service{}

	err = json.Unmarshal([]byte(content), &service)
	if err != nil {
		return errors.New(fmt.Sprintf("反序列化失败: %v\n", err))
	}

	_, err = client.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("更新Service失败: %v", err))
	}

	return nil
}

func (s *services) CreateService(client *kubernetes.Clientset, serviceCreate *ServiceCreate) (err error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceCreate.Name,
			Namespace: serviceCreate.Namespace,
			Labels:    serviceCreate.Label,
		},
		Spec: corev1.ServiceSpec{
			Selector: serviceCreate.Selector,
			Type:     corev1.ServiceType(serviceCreate.Type),
			Ports: []corev1.ServicePort{
				{
					Name:       serviceCreate.PortsName,
					Port:       serviceCreate.PortsPort,
					Protocol:   corev1.Protocol(serviceCreate.Protocol),
					TargetPort: intstr.FromInt(int(serviceCreate.TargetPort)),
				},
			},
		},
		Status: corev1.ServiceStatus{},
	}

	_, err = client.CoreV1().Services(serviceCreate.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("创建Service失败: %v", err))
	}

	return nil
}

func (s *services) DeleteService(client *kubernetes.Clientset, serviceName, namespace string) (err error) {
	if err = client.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{}); err != nil {
		return errors.New(fmt.Sprintf("删除Service失败: %v", err))
	}
	return nil
}

func (s *services) toCells(std []corev1.Service) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.ServiceCell(std[i])
	}
	return cells
}

func (s *services) fromCells(cells []dataselect.DataCell) []corev1.Service {
	services := make([]corev1.Service, len(cells))
	for i := range cells {
		services[i] = corev1.Service(cells[i].(dataselect.ServiceCell))
	}
	return services
}
