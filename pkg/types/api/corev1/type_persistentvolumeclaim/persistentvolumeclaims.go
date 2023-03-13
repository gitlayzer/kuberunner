package type_persistentvolumeclaim

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/utils/dataselect"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var PersistentVolumeClaim persistentvolumeclaim

type persistentvolumeclaim struct{}

func (p *persistentvolumeclaim) GetPersistentVolumeClaims(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (persistentvolumeclaimResp *PersistentVolumeClaimResp, err error) {
	persistentvolumeclaimList, err := client.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Pvc列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: p.toCells(persistentvolumeclaimList.Items),
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

	persistentvolumeclaims := p.fromCells(filtered.GenericDataList)

	return &PersistentVolumeClaimResp{
		Items: persistentvolumeclaims,
		Total: total,
	}, nil
}

func (p *persistentvolumeclaim) GetPersistentVolumeClaimDetail(client *kubernetes.Clientset, persistentvolumeclaimName, namespace string) (persistentvolumeclaim *corev1.PersistentVolumeClaim, err error) {
	persistentvolumeclaim, err = client.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), persistentvolumeclaimName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Pvc详情失败: %v", err))
	}

	return persistentvolumeclaim, nil
}

func (p *persistentvolumeclaim) UpdatePersistentVolumeClaim(client *kubernetes.Clientset, namespace, content string) (err error) {

	var persistentvolumeclaims = &corev1.PersistentVolumeClaim{}

	err = json.Unmarshal([]byte(content), &persistentvolumeclaims)
	if err != nil {
		return errors.New(fmt.Sprintf("反序列化失败: %v\n", err))
	}

	_, err = client.CoreV1().PersistentVolumeClaims(namespace).Update(context.TODO(), persistentvolumeclaims, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("更新Pvc失败: %v", err))
	}

	return nil
}

func (p *persistentvolumeclaim) CreatePersistentVolumeClaim(client *kubernetes.Clientset, persistentvolumeclaimCreate *PersistentVolumeClaimCreate) (err error) {
	persistentvolumeclaimcreate := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      persistentvolumeclaimCreate.Name,
			Namespace: persistentvolumeclaimCreate.Namespace,
			Labels:    persistentvolumeclaimCreate.Label,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.PersistentVolumeAccessMode(persistentvolumeclaimCreate.AccessMode),
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(persistentvolumeclaimCreate.Storage),
				},
			},
			StorageClassName: &persistentvolumeclaimCreate.StorageClass,
			Selector: &metav1.LabelSelector{
				MatchLabels: persistentvolumeclaimCreate.BindingPersistentVolume,
			},
		},
	}

	_, err = client.CoreV1().PersistentVolumeClaims(persistentvolumeclaimCreate.Namespace).Create(context.TODO(), persistentvolumeclaimcreate, metav1.CreateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("创建Pvc失败: %v", err))
	}

	return nil
}

func (p *persistentvolumeclaim) DeletePersistentVolumeClaim(client *kubernetes.Clientset, persistentvolumeclaimName, namespace string) (err error) {
	if err := client.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), persistentvolumeclaimName, metav1.DeleteOptions{}); err != nil {
		return errors.New(fmt.Sprintf("删除Pvc失败: %v", err))
	}

	return nil
}

func (p *persistentvolumeclaim) toCells(std []corev1.PersistentVolumeClaim) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.PersistentVolumeClaimCell(std[i])
	}
	return cells
}

func (p *persistentvolumeclaim) fromCells(cells []dataselect.DataCell) []corev1.PersistentVolumeClaim {
	persistentvolumeclaims := make([]corev1.PersistentVolumeClaim, len(cells))
	for i := range cells {
		persistentvolumeclaims[i] = corev1.PersistentVolumeClaim(cells[i].(dataselect.PersistentVolumeClaimCell))
	}
	return persistentvolumeclaims
}
