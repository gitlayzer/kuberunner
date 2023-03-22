package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/kuberunner/pkg/config"
	"github.com/gitlayzer/kuberunner/pkg/handlers/api"
	v1 "github.com/gitlayzer/kuberunner/pkg/handlers/api/v1"
	"github.com/gitlayzer/kuberunner/pkg/middlewares/cors"
	"github.com/gitlayzer/kuberunner/pkg/types/terminal"
	"github.com/gitlayzer/kuberunner/pkg/utils"
	"net/http"
)

func InitRouter() {
	r := gin.Default()

	r.Use(cors.Cors())

	r.POST("/login", api.Login.AuthLogin)

	r.GET("/api/v1/k8s/resources", v1.ResourceNum.GetResources)

	r.GET("/api/v1/k8s/cluster/list", v1.Cluster.GetClusterList)

	r.GET("/api/v1/k8s/event/list", v1.Event.GetEventList)

	r.GET("/api/v1/k8s/pod/list", v1.Pod.GetPodList)
	r.GET("/api/v1/k8s/pod/detail", v1.Pod.GetPodDetail)
	r.GET("/api/v1/k8s/pod/log", v1.Pod.GetPodLog)
	r.GET("/api/v1/k8s/pod/container", v1.Pod.GetPodContainer)
	r.PUT("/api/v1/k8s/pod/update", v1.Pod.UpdatePod)
	r.DELETE("/api/v1/k8s/pod/delete", v1.Pod.DeletePod)

	r.GET("/api/v1/k8s/deployment/list", v1.Deployment.GetDeployments)
	r.GET("/api/v1/k8s/deployment/detail", v1.Deployment.GetDeploymentDetail)
	r.PUT("/api/v1/k8s/deployment/update", v1.Deployment.UpdateDeployment)
	r.PUT("/api/v1/k8s/deployment/replicas", v1.Deployment.UpdateDeploymentReplicas)
	r.PUT("/api/v1/k8s/deployment/restart", v1.Deployment.RestartDeployment)
	r.POST("/api/v1/k8s/deployment/create", v1.Deployment.CreateDeployment)
	r.DELETE("/api/v1/k8s/deployment/delete", v1.Deployment.DeleteDeployment)

	r.GET("/api/v1/k8s/daemonset/list", v1.DaemonSet.GetDaemonSets)
	r.GET("/api/v1/k8s/daemonset/detail", v1.DaemonSet.GetDaemonSetDetail)
	r.PUT("/api/v1/k8s/daemonset/update", v1.DaemonSet.UpdateDaemonSet)
	r.PUT("/api/v1/k8s/daemonset/restart", v1.DaemonSet.RestartDaemonSet)
	r.POST("/api/v1/k8s/daemonset/create", v1.DaemonSet.CreateDaemonSet)
	r.DELETE("/api/v1/k8s/daemonset/delete", v1.DaemonSet.DeleteDaemonSet)

	r.GET("/api/v1/k8s/statefulset/list", v1.StatefulSet.GetStatefulSets)
	r.GET("/api/v1/k8s/statefulset/detail", v1.StatefulSet.GetStatefulSetDetail)
	r.PUT("/api/v1/k8s/statefulset/update", v1.StatefulSet.UpdateStatefulSet)
	r.PUT("/api/v1/k8s/statefulset/replicas", v1.StatefulSet.UpdateStatefulSetReplicas)
	r.PUT("/api/v1/k8s/statefulset/restart", v1.StatefulSet.RestartStatefulSet)
	r.POST("/api/v1/k8s/statefulset/create", v1.StatefulSet.CreateStatefulSet)
	r.DELETE("/api/v1/k8s/statefulset/delete", v1.StatefulSet.DeleteStatefulSet)

	r.GET("/api/v1/k8s/service/list", v1.Service.GetServices)
	r.GET("/api/v1/k8s/service/detail", v1.Service.GetServiceDetail)
	r.PUT("/api/v1/k8s/service/update", v1.Service.UpdateService)
	r.POST("/api/v1/k8s/service/create", v1.Service.CreateService)
	r.DELETE("/api/v1/k8s/service/delete", v1.Service.DeleteService)

	r.GET("/api/v1/k8s/ingress/list", v1.Ingress.GetIngresses)
	r.GET("/api/v1/k8s/ingress/detail", v1.Ingress.GetIngressDetail)
	r.PUT("/api/v1/k8s/ingress/update", v1.Ingress.UpdateIngress)
	r.POST("/api/v1/k8s/ingress/create", v1.Ingress.CreateIngress)
	r.DELETE("/api/v1/k8s/ingress/delete", v1.Ingress.DeleteIngress)

	r.GET("/api/v1/k8s/namespace/list", v1.Namespace.GetNamespaces)
	r.POST("/api/v1/k8s/namespace/create", v1.Namespace.CreateNamespace)
	r.DELETE("/api/v1/k8s/namespace/delete", v1.Namespace.DeleteNamespace)

	r.GET("/api/v1/k8s/node/list", v1.Node.GetNodes)

	r.GET("/api/v1/k8s/configmap/list", v1.ConfigMap.GetConfigMaps)
	r.GET("/api/v1/k8s/configmap/detail", v1.ConfigMap.GetConfigMapDetail)
	r.PUT("/api/v1/k8s/configmap/update", v1.ConfigMap.UpdateConfigMap)
	r.POST("/api/v1/k8s/configmap/create", v1.ConfigMap.CreateConfigMap)
	r.DELETE("/api/v1/k8s/configmap/delete", v1.ConfigMap.DeleteConfigMap)

	r.GET("/api/v1/k8s/secret/list", v1.Secret.GetSecrets)
	r.GET("/api/v1/k8s/secret/detail", v1.Secret.GetSecretDetail)
	r.PUT("/api/v1/k8s/secret/update", v1.Secret.UpdateSecret)
	r.POST("/api/v1/k8s/secret/create", v1.Secret.CreateSecret)
	r.DELETE("/api/v1/k8s/secret/delete", v1.Secret.DeleteSecret)

	r.GET("/api/v1/k8s/persistentvolume/list", v1.PersistentVolume.GetPersistentVolumes)
	r.GET("/api/v1/k8s/persistentvolume/detail", v1.PersistentVolume.GetPersistentVolumeDetail)
	r.PUT("/api/v1/k8s/persistentvolume/update", v1.PersistentVolume.UpdatePersistentVolume)
	r.POST("/api/v1/k8s/persistentvolume/create", v1.PersistentVolume.CreatePersistentVolume)
	r.DELETE("/api/v1/k8s/persistentvolume/delete", v1.PersistentVolume.DeletePersistentVolume)

	r.GET("/api/v1/k8s/persistentvolumeclaim/list", v1.PersistentVolumeClaim.GetPersistentVolumeClaims)
	r.GET("/api/v1/k8s/persistentvolumeclaim/detail", v1.PersistentVolumeClaim.GetPersistentVolumeClaimDetail)
	r.PUT("/api/v1/k8s/persistentvolumeclaim/update", v1.PersistentVolumeClaim.UpdatePersistentVolumeClaim)
	r.POST("/api/v1/k8s/persistentvolumeclaim/create", v1.PersistentVolumeClaim.CreatePersistentVolumeClaim)
	r.DELETE("/api/v1/k8s/persistentvolumeclaim/delete", v1.PersistentVolumeClaim.DeletePersistentVolumeClaim)

	r.GET("/api/v1/k8s/storageclass/list", v1.StorageClass.GetStorageClasses)
	r.GET("/api/v1/k8s/storageclass/detail", v1.StorageClass.GetStorageClassDetail)

	utils.Logo()

	wsHandler := http.NewServeMux()
	wsHandler.HandleFunc("/ws", terminal.Terminal.WsHandler)
	ws := &http.Server{
		Addr:    config.GetWsListenAddress(),
		Handler: wsHandler,
	}

	go func() {
		err := ws.ListenAndServe()
		if err != nil {
			return
		}
	}()

	go func() {
		err := r.Run(config.GetListenAddress())
		if err != nil {
			return
		}
	}()

	select {}
}
