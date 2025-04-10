package middlewares

import (
	"fast-gin/service/svc_redis"
	"fast-gin/utils/jwts"
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("token")
	_, err := jwts.ValidateJWT(token)
	if err != nil {
		response.FailWithMsg(c, "Authentication failed")
		c.Abort()
		return
	}
	if svc_redis.HasLoggedOut(token) {
		response.FailWithMsg(c, "User has logged out")
		c.Abort()
		return
	}
	c.Next()
}

func AdminAuthMiddleware(c *gin.Context) {
	token := c.GetHeader("token")

	claims, err := jwts.ValidateJWT(token)
	if err != nil {
		response.FailWithMsg(c, "Authentication failed")
		c.Abort()
		return
	}
	if svc_redis.HasLoggedOut(token) {
		response.FailWithMsg(c, "User has logged out")
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

func GetClaimsFrom(c *gin.Context) (claims *jwts.CustomClaims) {
	claims = new(jwts.CustomClaims)

	_claims, ok := c.Get("claims")
	if !ok {
		return
	}

	claims, ok = _claims.(*jwts.CustomClaims)
	if !ok {
		return
	}
	return
}
