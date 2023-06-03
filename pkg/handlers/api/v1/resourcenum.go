package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/kuberunner/pkg/types/api/corev1/type_resourcenum"
	"github.com/gitlayzer/kuberunner/pkg/utils"
	"net/http"
)

var ResourceNum resourcenum

type resourcenum struct{}

func (r *resourcenum) Register(router *gin.Engine) {
	resource := router.Group("/api/v1/k8s")
	{
		resource.GET("/resources", r.GetResources)
	}
}

func (r *resourcenum) GetResources(ctx *gin.Context) {
	params := new(struct {
		Cluster string `json:"cluster" form:"cluster"`
	})

	if err := ctx.Bind(params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	client, err := utils.K8s.GetClient(params.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	data, errs := type_resourcenum.ResourceNum.GetResourceNum(client)
	if len(errs) > 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": errs,
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    data,
	})
}
