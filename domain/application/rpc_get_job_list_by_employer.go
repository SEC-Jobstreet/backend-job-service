package application

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func (server *Server) GetJobListByEmployer(ctx context.Context, request *pb.JobListByEmployerRequest) (*pb.JobListResponse, error) {
	getJobListByEmployerAggregate, err := aggregate.GetJobListByEmployerAggregateGRPC(ctx, server.config, request)
	if err != nil {
		return nil, err
	}

	rsp, err := server.repo.GetJobListByEmployer(*getJobListByEmployerAggregate)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
