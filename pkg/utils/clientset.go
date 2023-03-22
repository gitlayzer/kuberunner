package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gitlayzer/kuberunner/pkg/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	ClientMap   map[string]*kubernetes.Clientset
	KubeConfMap map[string]string
}

func (k *k8s) GetClient(cluster string) (*kubernetes.Clientset, error) {
	client, ok := k.ClientMap[cluster]
	if !ok {
		return nil, errors.New(fmt.Sprintf("集群%s不存在，无法获取client\n", cluster))
	}
	return client, nil
}

func (k *k8s) Init() {
	mp := make(map[string]string, 0)

	k.ClientMap = make(map[string]*kubernetes.Clientset, 0)

	data, err := json.Marshal(config.GetKubeConfig())
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &mp)

	k.KubeConfMap = mp

	for key, value := range mp {
		conf, err := clientcmd.BuildConfigFromFlags("", value)

		if err != nil {
			panic(fmt.Sprintf("集群%s：创建K8s配置失败 %v\n", key, err))
		}

		clientSet, err := kubernetes.NewForConfig(conf)

		if err != nil {
			panic(fmt.Sprintf("集群%s：创建K8s客户端失败 %v\n", key, err))
		}

		k.ClientMap[key] = clientSet

		fmt.Printf("集群%s：创建K8s客户端成功\n", key)
	}
}
