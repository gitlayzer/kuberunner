package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/kuberunner/pkg/config"
	v1 "github.com/gitlayzer/kuberunner/pkg/handlers/api/v1"
	"github.com/gitlayzer/kuberunner/pkg/utils"
)

func InitRouter() {
	r := gin.Default()

	//r.Use(auth.Auth(), cors.Cors())

	r.GET("/api/v1/k8s/pod/list", v1.Pod.GetPodList)
	r.GET("/api/v1/k8s/pod/detail", v1.Pod.GetPodDetail)
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

	utils.Logo()

	err := r.Run(config.GetListenAddress())
	if err != nil {
		return
	}
}
