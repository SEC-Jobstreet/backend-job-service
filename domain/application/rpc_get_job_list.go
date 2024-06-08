package application

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func (server *Server) GetJobList(ctx context.Context, request *pb.JobListRequest) (*pb.JobListResponse, error) {
	getJobListAggregate, err := aggregate.GetJobListAggregateGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	rsp, err := server.repo.GetJobList(*getJobListAggregate)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
