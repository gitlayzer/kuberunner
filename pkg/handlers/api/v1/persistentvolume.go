package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/kuberunner/pkg/types/api/corev1/type_persistentvolume"
	"github.com/gitlayzer/kuberunner/pkg/utils"
	"net/http"
)

var PersistentVolume persistentvolume

type persistentvolume struct{}

func (p *persistentvolume) GetPersistentVolumes(ctx *gin.Context) {
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

	persistentvolumeResp, err := type_persistentvolume.PersistentVolume.GetPersistentVolumes(client, params.FilterName, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    persistentvolumeResp,
	})
}

func (p *persistentvolume) GetPersistentVolumeDetail(ctx *gin.Context) {
	params := new(struct {
		PersistentVolumeName string `form:"persistentvolume_name" json:"persistentvolume_name"`
		Cluster              string `form:"cluster" json:"cluster"`
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

	persistentvolumeResp, err := type_persistentvolume.PersistentVolume.GetPersistentVolumeDetail(client, params.PersistentVolumeName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    persistentvolumeResp,
	})
}

func (p *persistentvolume) UpdatePersistentVolume(ctx *gin.Context) {
	params := new(struct {
		Content string `json:"content"`
		Cluster string `json:"cluster"`
	})

	if err := ctx.ShouldBind(params); err != nil {
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

	if err := type_persistentvolume.PersistentVolume.UpdatePersistentVolume(client, params.Content); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    nil,
	})
}

func (p *persistentvolume) CreatePersistentVolume(ctx *gin.Context) {
	var (
		persistentvolumeCreate = new(type_persistentvolume.PersistentVolumeCreate)
		err                    error
	)

	if err := ctx.ShouldBindJSON(persistentvolumeCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	client, err := utils.K8s.GetClient(persistentvolumeCreate.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if err := type_persistentvolume.PersistentVolume.CreatePersistentVolume(client, persistentvolumeCreate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    nil,
	})
}

func (p *persistentvolume) DeletePersistentVolume(ctx *gin.Context) {
	params := new(struct {
		PersistentVolumeName string `json:"persistentvolume_name"`
		Cluster              string `json:"cluster"`
	})

	if err := ctx.ShouldBind(params); err != nil {
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

	if err := type_persistentvolume.PersistentVolume.DeletePersistentVolume(client, params.PersistentVolumeName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    nil,
	})
}
