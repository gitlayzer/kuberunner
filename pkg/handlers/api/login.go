package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gitlayzer/kuberunner/pkg/types/api/corev1/type_login"
	"net/http"
)

var Login login

type login struct{}

func (l *login) Register(router *gin.Engine) {
	router.POST("login", l.AuthLogin)
}

func (l *login) AuthLogin(ctx *gin.Context) {
	params := new(struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	})

	if err := ctx.ShouldBindJSON(params); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	err := type_login.Login.Login(params.Username, params.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "ok",
		"data": nil,
	})
}
