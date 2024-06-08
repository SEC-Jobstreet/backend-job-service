package application

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func (server *Server) GetJobListByAdmin(ctx context.Context, request *pb.JobListByAdminRequest) (*pb.JobListResponse, error) {
	getJobListByAdminAggregate, err := aggregate.GetJobListByAdminAggregateGRPC(ctx, server.config, request)
	if err != nil {
		return nil, err
	}

	rsp, err := server.repo.GetJobListByAdmin(*getJobListByAdminAggregate)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
