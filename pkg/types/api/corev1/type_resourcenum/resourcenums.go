package type_resourcenum

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var ResourceNum resourcenum

type resourcenum struct{}

func (r *resourcenum) GetResourceNum(client *kubernetes.Clientset) (map[string]int, []error) {
	var wg sync.WaitGroup
	wg.Add(12)

	errs := make([]error, 0)
	data := make(map[string]int, 0)
	go func() {
		list, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "Nodes", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "Pods", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "Services", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "Ingress", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "PersistentVolumes", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "PersistentVolumeClaims", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "ConfigMaps", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "Secrets", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "DaemonSets", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "StatefulSets", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "StorageClasses", len(list.Items))
		}
		wg.Done()
	}()
	go func() {
		list, err := client.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			errs = append(errs, err)
		} else {
			addMap(data, "Deployments", len(list.Items))
		}
		wg.Done()
	}()
	wg.Wait()
	return data, errs
}

func addMap(mp map[string]int, resource string, num int) {
	mt.Lock()
	defer mt.Unlock()
	mp[resource] = num
}
