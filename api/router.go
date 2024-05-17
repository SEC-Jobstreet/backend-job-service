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

	authRoutes.GET("/job/:id", s.GetJob)
	authRoutes.GET("/jobs", s.JobList)

	// GRPC
	// Employer Cognito
	authRoutes.POST("/post_job", middleware.AuthMiddleware(s.config, []string{"employers"}), s.PostJob)
	authRoutes.POST("/job_status", middleware.AuthMiddleware(s.config, []string{"employer", "admin"}), s.example)
	authRoutes.GET("/jobs_by_employer", middleware.AuthMiddleware(s.config, []string{"employer", "admin"}), s.GetJobByEmployer)

	s.router = router
}
