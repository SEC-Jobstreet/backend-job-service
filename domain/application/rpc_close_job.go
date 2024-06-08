package application

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func (server *Server) CloseJob(ctx context.Context, request *pb.CloseJobRequest) (*pb.CloseJobResponse, error) {
	closeJobAggregate, err := aggregate.CloseJobAggregateGRPC(ctx, server.config, request)
	if err != nil {
		return nil, err
	}

	rsp, err := server.repo.CloseJob(*closeJobAggregate)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
