package type_secret

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

var Secret secret

type secret struct{}

func (s *secret) GetSecrets(client *kubernetes.Clientset, filterName, namespace string, limit, page int) (secretResp *SecretResp, err error) {
	secretList, err := client.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Secret列表失败: %v", err))
	}

	selectableData := &dataselect.DataSelector{
		GenericDataList: s.toCells(secretList.Items),
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

	secrets := s.fromCells(filtered.GenericDataList)

	return &SecretResp{
		Items: secrets,
		Total: total,
	}, nil
}

func (s *secret) GetSecretDetail(client *kubernetes.Clientset, secretName, namespace string) (service *corev1.Secret, err error) {
	secret, err := client.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取Secret详情失败: %v", err))
	}

	return secret, nil
}

func (s *secret) UpdateSecret(client *kubernetes.Clientset, namespace, content string) (err error) {

	var secret = &corev1.Secret{}

	err = json.Unmarshal([]byte(content), &secret)
	if err != nil {
		return errors.New(fmt.Sprintf("反序列化失败: %v\n", err))
	}

	_, err = client.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("更新Secret失败: %v", err))
	}

	return nil
}

func (s *secret) CreateSecret(client *kubernetes.Clientset, secretCreate *SecretCreate) (err error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretCreate.Name,
			Namespace: secretCreate.Namespace,
			Labels:    secretCreate.Label,
		},
		Data: secretCreate.Data,
		Type: secretCreate.Type,
	}

	_, err = client.CoreV1().Secrets(secretCreate.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return errors.New(fmt.Sprintf("创建Secret失败: %v", err))
	}

	return nil
}

func (s *secret) DeleteSecret(client *kubernetes.Clientset, secretName, namespace string) (err error) {
	if err := client.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{}); err != nil {
		return errors.New(fmt.Sprintf("删除Secret失败: %v", err))
	}
	return nil
}

func (s *secret) toCells(std []corev1.Secret) []dataselect.DataCell {
	cells := make([]dataselect.DataCell, len(std))
	for i := range std {
		cells[i] = dataselect.SecretCell(std[i])
	}
	return cells
}

func (s *secret) fromCells(cells []dataselect.DataCell) []corev1.Secret {
	secrets := make([]corev1.Secret, len(cells))
	for i := range cells {
		secrets[i] = corev1.Secret(cells[i].(dataselect.SecretCell))
	}
	return secrets
}
