package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func LimitMiddleware(limit int32) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Enter LimitMiddleware")
		c.Next()
		fmt.Println("Exit LimitMiddleware")
	}
}
