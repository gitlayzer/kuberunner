package type_daemonset

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/utils/dataselect"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"time"
)

var DaemonSet daemonSet

type daemonSet struct{}

func (d *daemonSet) GetDaemonSets(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (daemonSetResp *DaemonSetResp, err error) {
	daemonSetList, err := client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取DaemonSet列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: d.toCells(daemonSetList.Items),
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

	daemonSets := d.fromCells(data.GenericDataList)

	return &DaemonSetResp{
		Items: daemonSets,
		Total: total,
	}, nil
}

func (d *daemonSet) GetDaemonSetDetail(client *kubernetes.Clientset, daemonSetName, namespace string) (daemonSet *appsv1.DaemonSet, err error) {
	daemonSet, err = client.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取DaemonSet详情失败: %v", err))
	}
	return daemonSet, nil
}

func (d *daemonSet) UpdateDaemonSet(client *kubernetes.Clientset, namespace, content string) (err error) {

	var daemonSet = &appsv1.DaemonSet{}

	err = json.Unmarshal([]byte(content), &daemonSet)
	if err != nil {
		return errors.New(fmt.Sprintf("反序列化失败: %v\n", err))
	}

	_, err = client.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("更新DaemonSet失败: %v", err))
	}

	return nil
}

func (d *daemonSet) RestartDaemonSet(client *kubernetes.Clientset, daemonSetName, namespace string) (err error) {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name": daemonSetName,
							"env": []map[string]string{{
								"name":  "RESTART_TIME",
								"value": time.Now().Format("2006-01-02 15:04:05"),
							}},
						},
					},
				},
			},
		},
	}

	patchByte, err := json.Marshal(patchData)
	if err != nil {
		return errors.New(fmt.Sprintf("序列化失败: %v\n", err))
	}

	_, err = client.AppsV1().DaemonSets(namespace).Patch(context.TODO(), daemonSetName, types.StrategicMergePatchType, patchByte, metav1.PatchOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("重启DaemonSet失败: %v", err))
	}

	return nil
}

func (d *daemonSet) CreateDaemonSet(client *kubernetes.Clientset, daemonSetCreate *DaemonSetCreate) (err error) {
	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      daemonSetCreate.Name,
			Namespace: daemonSetCreate.Namespace,
			Labels:    daemonSetCreate.Label,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: daemonSetCreate.Label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: daemonSetCreate.Label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            daemonSetCreate.Name,
							Image:           daemonSetCreate.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: daemonSetCreate.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		Status: appsv1.DaemonSetStatus{},
	}

	if daemonSetCreate.HealthCheck {
		daemonSet.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: daemonSetCreate.HealthPath,
					Port: intstr.FromInt(int(daemonSetCreate.ContainerPort)),
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      1,
			PeriodSeconds:       10,
		}
		daemonSet.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: daemonSetCreate.HealthPath,
					Port: intstr.FromInt(int(daemonSetCreate.ContainerPort)),
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      15,
			PeriodSeconds:       30,
		}
	}

	daemonSet.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(daemonSetCreate.Cpu),
		corev1.ResourceMemory: resource.MustParse(daemonSetCreate.Memory),
	}

	daemonSet.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(daemonSetCreate.Cpu),
		corev1.ResourceMemory: resource.MustParse(daemonSetCreate.Memory),
	}

	_, err = client.AppsV1().DaemonSets(daemonSetCreate.Namespace).Create(context.TODO(), daemonSet, metav1.CreateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("创建DaemonSet失败: %v", err))
	}

	return nil
}

func (d *daemonSet) DeleteDaemonSet(client *kubernetes.Clientset, daemonSetName, namespace string) (err error) {
	if err = client.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonSetName, metav1.DeleteOptions{}); err != nil {
		return errors.New(fmt.Sprintf("删除DaemonSet失败: %v", err))
	}
	return nil
}

func (d *daemonSet) toCells(std []appsv1.DaemonSet) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.DaemonSetCell(std[i])
	}
	return cells
}

func (d *daemonSet) fromCells(cells []dataselect.DataCell) []appsv1.DaemonSet {
	daemonsets := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		daemonsets[i] = appsv1.DaemonSet(cells[i].(dataselect.DaemonSetCell))
	}
	return daemonsets
}
