package api

import "github.com/gin-gonic/gin"

func InitApi() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.POST("/pipe", startPipe)
	router.GET("/pipe", getPipe)

	router.Run(":8080")
}
