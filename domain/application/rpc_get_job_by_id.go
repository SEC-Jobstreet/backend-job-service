package application

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func (server *Server) GetJobByID(ctx context.Context, request *pb.GetJobByIDRequest) (*pb.GetJobByIDResponse, error) {
	getJobByIDAggregate, err := aggregate.GetJobByIDAggregateGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	rsp, err := server.repo.GetJobByID(*getJobByIDAggregate)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
