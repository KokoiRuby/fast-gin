package middlewares

import (
	"fast-gin/utils/response"
	"github.com/gin-gonic/gin"
)

func BindJsonMiddleware[T any](c *gin.Context) {
	var req T
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithErr(c, err)
		c.Abort()
		return
	}
	// Set in context
	c.Set("request", req)
	return
}

func BindQueryMiddleware[T any](c *gin.Context) {
	var req T
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithErr(c, err)
		c.Abort()
		return
	}
	c.Set("request", req)
	return
}
func BindUriMiddleware[T any](c *gin.Context) {
	var req T
	err := c.ShouldBindUri(&req)
	if err != nil {
		response.FailWithErr(c, err)
		c.Abort()
		return
	}
	c.Set("request", req)
	return
}

func GetBind[T any](c *gin.Context) (cr T) {
	return c.MustGet("request").(T)
}
