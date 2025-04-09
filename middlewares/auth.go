package middlewares

import (
	"fast-gin/utils/jwt"
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	_, err := jwt.ValidateJWT(token)
	if err != nil {
		response.FailWithMsg(c, "Authentication failed")
		c.Abort()
		return
	}
	c.Next()
}

func AuthMiddlewareWithRole(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwt.ValidateJWT(token)
	if err != nil {
		response.FailWithMsg(c, "Authentication failed")
		c.Abort()
		return
	}
	if claims.RoleID != 1 {
		response.FailWithMsg(c, "Role Authentication failed")
		c.Abort()
		return
	}
	c.Next()
}
