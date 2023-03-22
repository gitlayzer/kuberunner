package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/kuberunner/pkg/types/api/corev1/type_persistentvolumeclaim"
	"github.com/gitlayzer/kuberunner/pkg/utils"
	"net/http"
)

var PersistentVolumeClaim persistentvolumeclaim

type persistentvolumeclaim struct{}

func (p *persistentvolumeclaim) GetPersistentVolumeClaims(ctx *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name" json:"filter_name"`
		Namespace  string `form:"namespace" json:"namespace"`
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

	persistentvolumeclaimResp, err := type_persistentvolumeclaim.PersistentVolumeClaim.GetPersistentVolumeClaims(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    persistentvolumeclaimResp,
	})
}

func (p *persistentvolumeclaim) GetPersistentVolumeClaimDetail(ctx *gin.Context) {
	params := new(struct {
		PersistentVolumeClaimName string `form:"persistent_volume_claim_name" json:"persistent_volume_claim_name"`
		Namespace                 string `form:"namespace" json:"namespace"`
		Cluster                   string `form:"cluster" json:"cluster"`
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

	persistentvolumeclaimResp, err := type_persistentvolumeclaim.PersistentVolumeClaim.GetPersistentVolumeClaimDetail(client, params.PersistentVolumeClaimName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    persistentvolumeclaimResp,
	})
}

func (p *persistentvolumeclaim) UpdatePersistentVolumeClaim(ctx *gin.Context) {
	params := new(struct {
		Namespace string `json:"namespace"`
		Content   string `json:"content"`
		Cluster   string `json:"cluster"`
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

	if err := type_persistentvolumeclaim.PersistentVolumeClaim.UpdatePersistentVolumeClaim(client, params.Namespace, params.Content); err != nil {
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

func (p *persistentvolumeclaim) CreatePersistentVolumeClaim(ctx *gin.Context) {
	var (
		persistentvolumeclaimCreate = new(type_persistentvolumeclaim.PersistentVolumeClaimCreate)
		err                         error
	)

	if err := ctx.ShouldBindJSON(persistentvolumeclaimCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	client, err := utils.K8s.GetClient(persistentvolumeclaimCreate.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if err := type_persistentvolumeclaim.PersistentVolumeClaim.CreatePersistentVolumeClaim(client, persistentvolumeclaimCreate); err != nil {
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

func (p *persistentvolumeclaim) DeletePersistentVolumeClaim(ctx *gin.Context) {
	params := new(struct {
		PersistentVolumeClaimName string `json:"persistent_volume_claim_name"`
		Namespace                 string `json:"namespace"`
		Cluster                   string `json:"cluster"`
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

	if err := type_persistentvolumeclaim.PersistentVolumeClaim.DeletePersistentVolumeClaim(client, params.PersistentVolumeClaimName, params.Namespace); err != nil {
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
