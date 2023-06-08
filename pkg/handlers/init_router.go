package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/kuberunner/pkg/config"
	"github.com/gitlayzer/kuberunner/pkg/handlers/api"
	v1 "github.com/gitlayzer/kuberunner/pkg/handlers/api/v1"
	"github.com/gitlayzer/kuberunner/pkg/middlewares/cors"
	"github.com/gitlayzer/kuberunner/pkg/types/terminal"
	"github.com/gitlayzer/kuberunner/pkg/utils"
)

func init() {
	utils.K8s.Init()
}

func InitRouter() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.Use(cors.Cors())

	api.Login.Register(r)
	v1.ResourceNum.Register(r)
	v1.Cluster.Register(r)
	v1.Event.Register(r)
	v1.Pod.Register(r)
	v1.Deployment.Register(r)
	v1.DaemonSet.Register(r)
	v1.StatefulSet.Register(r)
	v1.Service.Register(r)
	v1.Ingress.Register(r)
	v1.ConfigMap.Register(r)
	v1.Namespace.Register(r)
	v1.Node.Register(r)
	v1.Secret.Register(r)
	v1.PersistentVolume.Register(r)
	v1.PersistentVolumeClaim.Register(r)
	v1.StorageClass.Register(r)

	utils.Logo()
	r.GET("/ws", terminal.Terminal.ServeWsTerminal)

	err := r.Run(config.GetListenAddress())
	if err != nil {
		return
	}

}
