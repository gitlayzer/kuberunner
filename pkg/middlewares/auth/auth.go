package auth

import "github.com/gin-gonic/gin"

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{
				"code": 401,
				"msg":  "Token不能为空",
			})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
