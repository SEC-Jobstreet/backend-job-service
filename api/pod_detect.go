package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (s *Server) PodDetect(ctx *gin.Context) {
	podName := os.Getenv("HOSTNAME")
	if podName == "" {
		podName = "Unknown"
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello from pod: " + podName,
	})
}

func (s *Server) HelloWorld(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Hello, World!",
	})
}
