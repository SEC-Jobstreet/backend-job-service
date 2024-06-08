package application

import (
	"github.com/SEC-Jobstreet/backend-job-service/domain/repository"
	"github.com/SEC-Jobstreet/backend-job-service/domain/utils"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

type Server struct {
	pb.UnimplementedJobServiceServer
	config utils.Config
	repo   repository.JobRepository
}

// NewServer creates a new gRPC server.
func NewServer(config utils.Config, repo repository.JobRepository) (*Server, error) {

	server := &Server{
		config: config,
		repo:   repo,
	}

	return server, nil
}
