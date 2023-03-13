package type_deployment

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
	"strconv"
	"time"
)

var Deployment deployment

type deployment struct{}

func (d *deployment) GetDeployments(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (deploymentResp *DeploymentResp, err error) {
	deploymentList, err := client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Deployment列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: d.toCells(deploymentList.Items),
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

	deployments := d.fromCells(data.GenericDataList)

	return &DeploymentResp{
		Items: deployments,
		Total: total,
	}, nil
}

func (d *deployment) GetDeploymentDetail(client *kubernetes.Clientset, deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Deployment详情失败: %v", err))
	}

	return deployment, nil
}

func (d *deployment) UpdateDeployment(client *kubernetes.Clientset, namespace, content string) (err error) {
	var deployment = &appsv1.Deployment{}

	err = json.Unmarshal([]byte(content), &deployment)
	if err != nil {
		return errors.New(fmt.Sprintf("反序列化失败: %v\n", err))
	}

	_, err = client.AppsV1().Deployments(namespace).Update(context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("更新Deployment失败: %v\n", err))
	}

	return nil
}

func (d *deployment) UpdateDeploymentReplicas(client *kubernetes.Clientset, deploymentName, namespace string, replicas int32) (replica int32, err error) {
	scale, err := client.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return 0, errors.New(fmt.Sprintf("获取Deployment副本数失败: %v", err))
	}

	scale.Spec.Replicas = replicas

	_, err = client.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		return 0, errors.New(fmt.Sprintf("更新Deployment副本数失败: %v", err))
	}

	return replicas, nil
}

func (d *deployment) RestartDeployment(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name": deploymentName,
							"env": []map[string]string{{
								"name":  "RESTARTED_AT",
								"value": strconv.FormatInt(time.Now().Unix(), 10),
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

	_, err = client.AppsV1().Deployments(namespace).Patch(context.TODO(), deploymentName, types.StrategicMergePatchType, patchByte, metav1.PatchOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("重启Deployment失败: %v", err))
	}

	return nil
}

func (d *deployment) CreateDeployment(client *kubernetes.Clientset, data *DeploymentCreate) (err error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: data.Label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            data.Name,
							Image:           data.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          data.Name,
									ContainerPort: data.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}

	if data.HealthCheck {
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.FromInt(int(data.ContainerPort)),
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      1,
			PeriodSeconds:       10,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.FromInt(int(data.ContainerPort)),
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      15,
			PeriodSeconds:       30,
		}
	}

	deployment.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}
	deployment.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	_, err = client.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("创建Deployment失败: %v", err))
	}

	return nil
}

func (d *deployment) DeleteDeployment(client *kubernetes.Clientset, deploymentName, namespace string) (err error) {
	if err = client.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{}); err != nil {
		return errors.New(fmt.Sprintf("删除Deployment失败: %v", err))
	}

	return nil
}

func (d *deployment) toCells(std []appsv1.Deployment) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.DeploymentCell(std[i])
	}
	return cells
}

func (d *deployment) fromCells(cells []dataselect.DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(dataselect.DeploymentCell))
	}
	return deployments
}
