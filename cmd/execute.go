package cmd

import (
	"github.com/gitlayzer/kuberunner/pkg/handlers"
	"github.com/gitlayzer/kuberunner/pkg/utils"
)

func Execute() {
	utils.K8s.Init()
	handlers.InitRouter()
}
