package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/kuberunner/pkg/types/api/corev1/type_namespace"
	"github.com/gitlayzer/kuberunner/pkg/utils"
	"net/http"
)

var Namespace namespace

type namespace struct{}

func (n *namespace) GetNamespaces(ctx *gin.Context) {
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

	namespaceResp, err := type_namespace.Namespace.GetNamespaces(client, params.FilterName, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    namespaceResp,
	})
}

func (n *namespace) CreateNamespace(ctx *gin.Context) {
	var (
		namespaceCreate = new(type_namespace.NamespaceCreate)
		err             error
	)

	if err := ctx.ShouldBindJSON(namespaceCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	client, err := utils.K8s.GetClient(namespaceCreate.Cluster)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if err := type_namespace.Namespace.CreateNamespace(client, namespaceCreate); err != nil {
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

func (n *namespace) DeleteNamespace(ctx *gin.Context) {
	params := new(struct {
		NamespaceName string `json:"namespace_name"`
		Cluster       string `json:"cluster"`
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

	if err := type_namespace.Namespace.DeleteNamespace(client, params.NamespaceName); err != nil {
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
