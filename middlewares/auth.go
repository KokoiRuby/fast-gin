package middlewares

import (
	"fast-gin/utils/jwts"
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	_, err := jwts.ValidateJWT(token)
	if err != nil {
		response.FailWithMsg(c, "Authentication failed")
		c.Abort()
		return
	}
	c.Next()
}

func AuthMiddlewareWithRole(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := jwts.ValidateJWT(token)
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

	// Set claim in context
	c.Set("claims", claims)
	c.Next()
}

func GetClaimsFrom(c *gin.Context) (claims *jwt.Claims) {
	claims = new(jwt.Claims)

	_claims, ok := c.Get("claims")
	if !ok {
		return
	}

	claims, ok = _claims.(*jwt.Claims)
	if !ok {
		return
	}
	return
}
