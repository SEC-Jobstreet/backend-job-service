package application

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func (server *Server) GetNumberOfJob(ctx context.Context, request *pb.GetNumberOfJobRequest) (*pb.GetNumberOfJobResponse, error) {
	rsp := server.repo.GetNumberOfJob()
	return rsp, nil
}
