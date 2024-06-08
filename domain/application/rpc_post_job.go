package application

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate"
	"github.com/SEC-Jobstreet/backend-job-service/domain/factory"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func (server *Server) PostJob(ctx context.Context, request *pb.PostJobRequest) (*pb.PostJobResponse, error) {
	postJobAggregate, err := aggregate.PostJobAggregateGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	if err = server.repo.Create(*postJobAggregate); err != nil {
		return nil, err
	}

	return &pb.PostJobResponse{
		Job: factory.ConvertJob(postJobAggregate.Job),
	}, nil
}
