package type_event

import (
	"context"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/utils/dataselect"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Event event

type event struct{}

func (e *event) GetEventList(client *kubernetes.Clientset, filterName string, limit, page int) (eventResp *EventResp, err error) {
	eventList, err := client.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Event列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: e.toCells(eventList.Items),
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
	events := e.fromCells(data.GenericDataList)

	return &EventResp{
		Items: events,
		Total: total,
	}, nil
}

func (e *event) toCells(std []corev1.Event) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.EventCell(std[i])
	}
	return cells
}

func (e *event) fromCells(cells []dataselect.DataCell) []corev1.Event {
	events := make([]corev1.Event, len(cells))
	for i := range cells {
		events[i] = corev1.Event(cells[i].(dataselect.EventCell))
	}
	return events
}
