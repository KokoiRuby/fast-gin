package routers

import (
	"fast-gin/global"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Run() {
	gin.SetMode(global.Config.Gin.Mode)

	r := gin.Default()

	// Static route
	// curl http://localhost:8080/uploads/test.txt
	r.Static("/uploads", "./static/uploads")

	// Grouping routes
	v1 := r.Group("v1")
	UserRouter(v1)
	ImageRouter(v1)

	// Run Gin server
	err := r.Run(global.Config.Gin.Addr())
	if err != nil {
		logrus.Fatalf("Failed to start Gin server: %v", err)
		return
	}
}
