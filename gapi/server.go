package gapi

import (
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"gorm.io/gorm"
)

type Server struct {
	pb.UnimplementedJobServiceServer
	config utils.Config
	store  *gorm.DB
}

// NewServer creates a new gRPC server.
func NewServer(config utils.Config, store *gorm.DB) (*Server, error) {

	server := &Server{
		config: config,
		store:  store,
	}

	return server, nil
}
