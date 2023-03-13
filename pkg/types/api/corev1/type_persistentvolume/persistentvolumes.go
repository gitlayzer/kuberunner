package type_persistentvolume

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

var PersistentVolume persistentvolume

type persistentvolume struct{}

func (p *persistentvolume) GetPersistentVolumes(client *kubernetes.Clientset, filterName string, limit, page int) (persistentvolumeResp *PersistentVolumeResp, err error) {
	persistentvolumeList, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Pv列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: p.toCells(persistentvolumeList.Items),
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

	persistentvolumes := p.fromCells(filtered.GenericDataList)

	return &PersistentVolumeResp{
		Items: persistentvolumes,
		Total: total,
	}, nil
}

func (p *persistentvolume) GetPersistentVolumeDetail(client *kubernetes.Clientset, pvName string) (service *corev1.PersistentVolume, err error) {
	persistentvolumes, err := client.CoreV1().PersistentVolumes().Get(context.TODO(), pvName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Pv详情失败: %v", err))
	}

	return persistentvolumes, nil
}

func (p *persistentvolume) UpdatePersistentVolume(client *kubernetes.Clientset, content string) (err error) {

	var persistentvolumes = &corev1.PersistentVolume{}

	err = json.Unmarshal([]byte(content), &persistentvolumes)
	if err != nil {
		return errors.New(fmt.Sprintf("反序列化失败: %v\n", err))
	}

	_, err = client.CoreV1().PersistentVolumes().Update(context.TODO(), persistentvolumes, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("更新Pv失败: %v", err))
	}

	return nil
}

func (p *persistentvolume) CreatePersistentVolume(client *kubernetes.Clientset, pvCreate *PersistentVolumeCreate) (err error) {
	persistentvolumeCrete := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name:   pvCreate.Name,
			Labels: pvCreate.Label,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName:              pvCreate.StorageClass,
			AccessModes:                   []corev1.PersistentVolumeAccessMode{corev1.PersistentVolumeAccessMode(pvCreate.AccessMode)},
			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimPolicy(pvCreate.PersistentVolumeReclaimPolicy),
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(pvCreate.Storage),
			},
		},
		Status: corev1.PersistentVolumeStatus{},
	}

	// 判断传递的类型
	if pvCreate.Nfs.Path != "" {
		persistentvolumeCrete.Spec.PersistentVolumeSource = corev1.PersistentVolumeSource{
			NFS: &corev1.NFSVolumeSource{
				Path:   pvCreate.Nfs.Path,
				Server: pvCreate.Nfs.Server,
			},
		}
	} else if pvCreate.Rbd.Image != "" {
		persistentvolumeCrete.Spec.PersistentVolumeSource = corev1.PersistentVolumeSource{
			RBD: &corev1.RBDPersistentVolumeSource{
				CephMonitors: pvCreate.Rbd.Monitors,
				RBDPool:      pvCreate.Rbd.Pool,
				RBDImage:     pvCreate.Rbd.Image,
				RadosUser:    pvCreate.Rbd.User,
				SecretRef:    &pvCreate.Rbd.SecretRef,
				ReadOnly:     pvCreate.Rbd.ReadOnly,
				FSType:       pvCreate.Rbd.FsType,
			},
		}
	} else if pvCreate.HostPath.Path != "" {
		persistentvolumeCrete.Spec.PersistentVolumeSource = corev1.PersistentVolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: pvCreate.HostPath.Path,
			},
		}
	} else if pvCreate.Local.Path != "" {
		persistentvolumeCrete.Spec.PersistentVolumeSource = corev1.PersistentVolumeSource{
			Local: &corev1.LocalVolumeSource{
				Path:   pvCreate.Local.Path,
				FSType: pvCreate.Local.FsType,
			},
		}
	} else if pvCreate.CephFs.Path != "" {
		persistentvolumeCrete.Spec.PersistentVolumeSource = corev1.PersistentVolumeSource{
			CephFS: &corev1.CephFSPersistentVolumeSource{
				Monitors:  pvCreate.CephFs.Monitors,
				Path:      pvCreate.CephFs.Path,
				User:      pvCreate.CephFs.User,
				SecretRef: &pvCreate.CephFs.SecretRef,
			},
		}
	}

	_, err = client.CoreV1().PersistentVolumes().Create(context.TODO(), persistentvolumeCrete, metav1.CreateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("创建Pv失败: %v", err))
	}

	return nil
}

func (p *persistentvolume) DeletePersistentVolume(client *kubernetes.Clientset, pvName string) (err error) {
	if err := client.CoreV1().PersistentVolumes().Delete(context.TODO(), pvName, metav1.DeleteOptions{}); err != nil {
		return errors.New(fmt.Sprintf("删除Pv失败: %v", err))
	}
	return nil
}

func (p *persistentvolume) toCells(std []corev1.PersistentVolume) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.PersistentVolumeCell(std[i])
	}
	return cells
}

func (p *persistentvolume) fromCells(cells []dataselect.DataCell) []corev1.PersistentVolume {
	persistentvolumes := make([]corev1.PersistentVolume, len(cells))
	for i := range cells {
		persistentvolumes[i] = corev1.PersistentVolume(cells[i].(dataselect.PersistentVolumeCell))
	}
	return persistentvolumes
}
