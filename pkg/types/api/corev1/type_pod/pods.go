package type_pod

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/utils/dataselect"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Pod pod

type pod struct{}

func (p *pod) GetPods(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	podList, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Pod列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: p.toCells(podList.Items),
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

	pods := p.fromCells(data.GenericDataList)

	return &PodsResp{
		Items: pods,
		Total: total,
	}, nil
}

func (p *pod) GetPodDetail(client *kubernetes.Clientset, podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Pod详情失败: %v", err))
	}

	return pod, nil
}

func (p *pod) GetPodContainer(client *kubernetes.Clientset, name, namespace string) (containers []string, err error) {
	pod, err := p.GetPodDetail(client, name, namespace)
	if err != nil {
		return nil, err
	}

	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}

	return containers, nil
}

func (p *pod) GetPodLog(client *kubernetes.Clientset, containerName, podName, namespace string) (log string, err error) {
	lineLimit := int64(1000)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimit,
	}

	req := client.CoreV1().Pods(namespace).GetLogs(podName, option)

	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "", errors.New("获取PodLog失败, " + err.Error())
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", errors.New("复制PodLog失败, " + err.Error())
	}

	return buf.String(), nil
}

func (p *pod) UpdatePod(client *kubernetes.Clientset, namespace, content string) (err error) {
	var pod corev1.Pod

	err = json.Unmarshal([]byte(content), &pod)
	if err != nil {
		return errors.New(fmt.Sprintf("序列化失败: %v", err))
	}

	_, err = client.CoreV1().Pods(namespace).Update(context.Background(), &pod, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("更新Pod失败: %v", err))
	}

	return nil
}

func (p *pod) DeletePod(client *kubernetes.Clientset, podName, namespace string) (err error) {
	if err = client.CoreV1().Pods(namespace).Delete(context.Background(), podName, metav1.DeleteOptions{}); err != nil {
		return errors.New(fmt.Sprintf("删除Pod失败: %v", err))
	}

	return nil
}

func (p *pod) toCells(std []corev1.Pod) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.PodCell(std[i])
	}
	return cells
}

func (p *pod) fromCells(cells []dataselect.DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(dataselect.PodCell))
	}
	return pods
}
