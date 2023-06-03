package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/kuberunner/pkg/types/api/corev1/type_statefulset"
	"github.com/gitlayzer/kuberunner/pkg/utils"
	"net/http"
)

var (
	StatefulSet statefulset
)

type statefulset struct{}

func (s *statefulset) Register(router *gin.Engine) {
	statefulset := router.Group("/api/v1/k8s/statefulset")
	{
		statefulset.GET("/list", s.GetStatefulSets)
		statefulset.GET("/detail", s.GetStatefulSetDetail)
		statefulset.PUT("/update", s.UpdateStatefulSet)
		statefulset.PUT("/replicas", s.UpdateStatefulSetReplicas)
		statefulset.PUT("/restart", s.RestartStatefulSet)
		statefulset.DELETE("/delete", s.DeleteStatefulSet)
		statefulset.POST("/create", s.CreateStatefulSet)
	}
}

func (s *statefulset) GetStatefulSets(ctx *gin.Context) {
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

	statefulSetResp, err := type_statefulset.StatefulSet.GetStatefulSets(client, params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    statefulSetResp,
	})
}

func (s *statefulset) GetStatefulSetDetail(ctx *gin.Context) {
	params := new(struct {
		StatefulSetName string `form:"statefulset_name" json:"statefulset_name"`
		Namespace       string `form:"namespace" json:"namespace"`
		Cluster         string `form:"cluster" json:"cluster"`
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

	statefulSetDetail, err := type_statefulset.StatefulSet.GetStatefulSetDetail(client, params.StatefulSetName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    statefulSetDetail,
	})
}

func (s *statefulset) UpdateStatefulSet(ctx *gin.Context) {
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

	if err := type_statefulset.StatefulSet.UpdateStatefulSet(client, params.Namespace, params.Content); err != nil {
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

func (s *statefulset) UpdateStatefulSetReplicas(ctx *gin.Context) {
	params := new(struct {
		StatefulSetName string `json:"statefulset_name"`
		Namespace       string `json:"namespace"`
		Cluster         string `json:"cluster"`
		Replicas        int32  `json:"replicas"`
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

	replicas, err := type_statefulset.StatefulSet.UpdateStatefulSetReplicas(client, params.StatefulSetName, params.Namespace, params.Replicas)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    replicas,
	})
}

func (s *statefulset) RestartStatefulSet(ctx *gin.Context) {
	params := new(struct {
		StatefulSetName string `json:"statefulset_name"`
		Namespace       string `json:"namespace"`
		Cluster         string `json:"cluster"`
	})

	if err := ctx.ShouldBindJSON(params); err != nil {
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

	if err := type_statefulset.StatefulSet.RestartStatefulSet(client, params.StatefulSetName, params.Namespace); err != nil {
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

func (d *statefulset) CreateStatefulSet(ctx *gin.Context) {
	var (
		statefulSetCreate = new(type_statefulset.StatefulSetCreate)
		err               error
	)

	if err := ctx.ShouldBindJSON(statefulSetCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	client, err := utils.K8s.GetClient(statefulSetCreate.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if err := type_statefulset.StatefulSet.CreateStatefulSet(client, statefulSetCreate); err != nil {
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

func (d *statefulset) DeleteStatefulSet(ctx *gin.Context) {
	params := new(struct {
		StatefulSetName string `json:"statefulset_name"`
		Namespace       string `json:"namespace"`
		Cluster         string `json:"cluster"`
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

	if err := type_statefulset.StatefulSet.DeleteStatefulSet(client, params.StatefulSetName, params.Namespace); err != nil {
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
