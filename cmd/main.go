package main

import (
	"github.com/gitlayzer/kuberunner/pkg/handlers"
)

func Execute() {
	handlers.InitRouter()
}

func main() {
	Execute()
}
