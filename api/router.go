package api

import (
	"time"

	"github.com/SEC-Jobstreet/backend-job-service/api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) setupRouter() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "content-type", "accept", "authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authRoutes := router.Group("/api/v1")

	authRoutes.GET("/job/:id", s.example)
	authRoutes.GET("/jobs", s.example)

	authRoutes.POST("/post_job", middleware.AuthMiddleware(s.config), s.PostJob)
	authRoutes.POST("/job_accepted", middleware.AuthMiddleware(s.config), s.example)
	authRoutes.POST("/job_closed", middleware.AuthMiddleware(s.config), s.example)
	authRoutes.GET("/job_by_employer_id/:id", middleware.AuthMiddleware(s.config), s.example)

	s.router = router
}
