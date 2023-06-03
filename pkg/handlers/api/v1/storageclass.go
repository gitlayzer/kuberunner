package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/kuberunner/pkg/types/api/corev1/type_storageclass"
	"github.com/gitlayzer/kuberunner/pkg/utils"
	"net/http"
)

var StorageClass storageClass

type storageClass struct{}

func (s *storageClass) Register(router *gin.Engine) {
	storageClass := router.Group("/api/v1/k8s/storageclass")
	{
		storageClass.GET("/list", s.GetStorageClasses)
		storageClass.GET("/detail", s.GetStorageClassDetail)
	}
}

func (s *storageClass) GetStorageClasses(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name" json:"filter_name"`
		Page       int    `form:"page" json:"page"`
		Limit      int    `form:"limit" json:"limit"`
		Cluster    string `form:"cluster" json:"cluster"`
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

	storageClassResp, err := type_storageclass.StorageClass.GetStorageClasses(client, params.FilterName, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    storageClassResp,
	})
}

func (s *storageClass) GetStorageClassDetail(ctx *gin.Context) {
	params := new(struct {
		StorageClassName string `form:"storageclass_name" json:"storageclass_name"`
		Cluster          string `form:"cluster" json:"cluster"`
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

	storageClassDetail, err := type_storageclass.StorageClass.GetStorageClassDetail(client, params.StorageClassName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    storageClassDetail,
	})
}
