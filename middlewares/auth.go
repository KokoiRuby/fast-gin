package middlewares

import (
	"fast-gin/utils/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	_, err := jwt.ValidateJWT(token)
	if err != nil {
		c.JSON(200, gin.H{"code": 7, "msg": "Authentication failed", "data": gin.H{}})
		c.Abort()
		return
	}
	c.Next()
}

func AuthMiddlewareWithRole(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwt.ValidateJWT(token)
	if err != nil {
		c.JSON(200, gin.H{"code": 7, "msg": "Authentication failed", "data": gin.H{}})
		c.Abort()
		return
	}
	if claims.RoleID != 1 {
		c.JSON(200, gin.H{"code": 7, "msg": "Admin Authentication failed", "data": gin.H{}})
		c.Abort()
		return
	}
	c.Next()
}
