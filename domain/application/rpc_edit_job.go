package application

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func (server *Server) EditJob(ctx context.Context, request *pb.EditJobRequest) (*pb.EditJobResponse, error) {
	editJobAggregate, err := aggregate.EditJobAggregateGRPC(ctx, server.config, request)
	if err != nil {
		return nil, err
	}

	rsp, err := server.repo.EditJob(*editJobAggregate)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
