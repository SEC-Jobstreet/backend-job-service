package application

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func (server *Server) GetNumberOfNewJob(ctx context.Context, request *pb.GetNumberOfNewJobRequest) (*pb.GetNumberOfNewJobResponse, error) {
	numberOfNewJobAggregate := aggregate.GetNumberOfNewJobAggregateGRPC(ctx, request)
	rsp := server.repo.GetNumberOfNewJob(*numberOfNewJobAggregate)
	return rsp, nil
}
