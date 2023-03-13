package type_storageclass

import (
	"context"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/utils/dataselect"
	v1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var StorageClass storageClass

type storageClass struct{}

func (s *storageClass) GetStorageClasses(client *kubernetes.Clientset, filterName string, limit, page int) (storageClassResp *StorageClassResp, err error) {
	storageClassList, err := client.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取StorageClass列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: s.toCells(storageClassList.Items),
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

	storageClasses := s.fromCells(data.GenericDataList)

	return &StorageClassResp{
		Items: storageClasses,
		Total: total,
	}, nil
}

func (s *storageClass) GetStorageClassDetail(client *kubernetes.Clientset, serviceName string) (storageClass *v1.StorageClass, err error) {
	storageclass, err := client.StorageV1().StorageClasses().Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取StorageClass详情失败: %v", err))
	}

	return storageclass, nil
}

func (s *storageClass) toCells(std []v1.StorageClass) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.StorageClassCell(std[i])
	}
	return cells
}

func (s *storageClass) fromCells(cells []dataselect.DataCell) []v1.StorageClass {
	storageClasses := make([]v1.StorageClass, len(cells))
	for i := range cells {
		storageClasses[i] = v1.StorageClass(cells[i].(dataselect.StorageClassCell))
	}
	return storageClasses
}
