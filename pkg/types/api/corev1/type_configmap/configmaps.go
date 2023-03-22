package type_configmap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/utils/dataselect"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Configmap configmap

type configmap struct{}

func (c *configmap) GetConfigMaps(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (configmapResp *ConfigMapResp, err error) {
	configmapList, err := client.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Configmap列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: c.toCells(configmapList.Items),
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

	configmaps := c.fromCells(data.GenericDataList)

	return &ConfigMapResp{
		Items: configmaps,
		Total: total,
	}, nil
}

func (c *configmap) GetConfigMapDetail(client *kubernetes.Clientset, configmapName, namespace string) (service *corev1.ConfigMap, err error) {
	configmap, err := client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), configmapName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取ConfigMap详情失败: %v", err))
	}

	return configmap, nil
}

func (c *configmap) UpdateConfigMap(client *kubernetes.Clientset, namespace, content string) (err error) {

	var configmap = &corev1.ConfigMap{}

	err = json.Unmarshal([]byte(content), &configmap)
	if err != nil {
		return errors.New(fmt.Sprintf("反序列化失败: %v\n", err))
	}

	_, err = client.CoreV1().ConfigMaps(namespace).Update(context.TODO(), configmap, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("更新ConfigMap失败: %v", err))
	}

	return nil
}

func (c *configmap) CreateConfigMap(client *kubernetes.Clientset, configmapCreate *ConfigMapCreate) (err error) {
	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configmapCreate.Name,
			Namespace: configmapCreate.Namespace,
		},
		Data: configmapCreate.Data,
	}

	_, err = client.CoreV1().ConfigMaps(configmapCreate.Namespace).Create(context.TODO(), configmap, metav1.CreateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("创建ConfigMap失败: %v", err))
	}

	return nil
}

func (c *configmap) DeleteConfigMap(client *kubernetes.Clientset, configmapName, namespace string) (err error) {
	if err := client.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), configmapName, metav1.DeleteOptions{}); err != nil {
		return errors.New(fmt.Sprintf("删除ConfigMap失败: %v", err))
	}

	return nil
}

func (c *configmap) toCells(std []corev1.ConfigMap) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.ConfigMapCell(std[i])
	}
	return cells
}

func (c *configmap) fromCells(cells []dataselect.DataCell) []corev1.ConfigMap {
	configmaps := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		configmaps[i] = corev1.ConfigMap(cells[i].(dataselect.ConfigMapCell))
	}
	return configmaps
}
