package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	fmt.Println("Enter AuthMiddleware")
	c.Next()
	fmt.Println("Exit AuthMiddleware")
}
