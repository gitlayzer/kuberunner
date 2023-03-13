package type_node

import (
	"context"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/utils/dataselect"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var Node node

type node struct{}

func (n *node) GetNodes(client *kubernetes.Clientset, filterName string, limit, page int) (nodeResp *NodeResp, err error) {
	nodeList, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Node列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: n.toCells(nodeList.Items),
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

	nodes := n.fromCells(filtered.GenericDataList)

	return &NodeResp{
		Items: nodes,
		Total: total,
	}, nil
}

func (n *node) toCells(std []corev1.Node) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.NodeCell(std[i])
	}
	return cells
}

func (n *node) fromCells(cells []dataselect.DataCell) []corev1.Node {
	nodes := make([]corev1.Node, len(cells))
	for i := range cells {
		nodes[i] = corev1.Node(cells[i].(dataselect.NodeCell))
	}
	return nodes
}
