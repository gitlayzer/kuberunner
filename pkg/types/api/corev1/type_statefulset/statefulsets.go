package type_statefulset

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

var StatefulSet statefulset

type statefulset struct{}

func (d *statefulset) GetStatefulSets(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (statefulSetResp *StatefulSetResp, err error) {
	statefulSetList, err := client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取DaemonSet列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: d.toCells(statefulSetList.Items),
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

	statefulSets := d.fromCells(data.GenericDataList)

	return &StatefulSetResp{
		Items: statefulSets,
		Total: total,
	}, nil
}

func (d *statefulset) GetStatefulSetDetail(client *kubernetes.Clientset, statefulSetName, namespace string) (statefulSet *appsv1.StatefulSet, err error) {
	statefulSet, err = client.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取StatefulSet详情失败: %v", err))
	}
	return statefulSet, nil
}

func (d *statefulset) UpdateStatefulSet(client *kubernetes.Clientset, namespace, content string) (err error) {

	var statefulSet = &appsv1.StatefulSet{}

	err = json.Unmarshal([]byte(content), &statefulSet)
	if err != nil {
		return errors.New(fmt.Sprintf("反序列化失败: %v\n", err))
	}

	_, err = client.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("更新DaemonSet失败: %v", err))
	}

	return nil
}

func (d *statefulset) UpdateStatefulSetReplicas(client *kubernetes.Clientset, statefulSetName, namespace string, replicas int32) (replica int32, err error) {
	scale, err := client.AppsV1().StatefulSets(namespace).GetScale(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		return 0, errors.New(fmt.Sprintf("获取Statefulset副本数失败: %v", err))
	}

	scale.Spec.Replicas = replicas

	_, err = client.AppsV1().StatefulSets(namespace).UpdateScale(context.TODO(), statefulSetName, scale, metav1.UpdateOptions{})
	if err != nil {
		return 0, errors.New(fmt.Sprintf("更新Deployment副本数失败: %v", err))
	}

	return replicas, nil
}

func (d *statefulset) RestartStatefulSet(client *kubernetes.Clientset, statefulSetName, namespace string) (err error) {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name": statefulSetName,
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

	_, err = client.AppsV1().StatefulSets(namespace).Patch(context.TODO(), statefulSetName, types.StrategicMergePatchType, patchByte, metav1.PatchOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("重启statefulSet失败: %v", err))
	}

	return nil
}

func (d *statefulset) CreateStatefulSet(client *kubernetes.Clientset, statefulset *StatefulSetCreate) (err error) {
	statefulet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulset.Name,
			Namespace: statefulset.Namespace,
			Labels:    statefulset.Label,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &statefulset.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: statefulset.Label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: statefulset.Label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            statefulset.Name,
							Image:           statefulset.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          statefulset.Name,
									ContainerPort: statefulset.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		Status: appsv1.StatefulSetStatus{},
	}

	if statefulset.HealthCheck {
		statefulet.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: statefulset.HealthPath,
					Port: intstr.FromInt(int(statefulset.ContainerPort)),
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      1,
			PeriodSeconds:       10,
		}
		statefulet.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: statefulset.HealthPath,
					Port: intstr.FromInt(int(statefulset.ContainerPort)),
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      15,
			PeriodSeconds:       30,
		}
	}

	statefulet.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(statefulset.Cpu),
		corev1.ResourceMemory: resource.MustParse(statefulset.Memory),
	}
	statefulet.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(statefulset.Cpu),
		corev1.ResourceMemory: resource.MustParse(statefulset.Memory),
	}

	_, err = client.AppsV1().StatefulSets(statefulset.Namespace).Create(context.TODO(), statefulet, metav1.CreateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("创建Statefulset失败: %v", err))
	}

	return nil
}

func (d *statefulset) DeleteStatefulSet(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	if err = client.AppsV1().StatefulSets(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{}); err != nil {
		return errors.New(fmt.Sprintf("删除Statefulset失败: %v", err))
	}

	return nil
}

func (d *statefulset) toCells(std []appsv1.StatefulSet) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.StatefulSetCell(std[i])
	}
	return cells
}

func (d *statefulset) fromCells(cells []dataselect.DataCell) []appsv1.StatefulSet {
	statefulsets := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		statefulsets[i] = appsv1.StatefulSet(cells[i].(dataselect.StatefulSetCell))
	}
	return statefulsets
}
