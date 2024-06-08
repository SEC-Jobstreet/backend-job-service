package application

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func (server *Server) ChangeStatusJobByAdmin(ctx context.Context, request *pb.ChangeStatusJobByAdminRequest) (*pb.ChangeStatusJobByAdminResponse, error) {
	changeStatusJobByAdminAggregate, err := aggregate.ChangeStatusJobByAdminAggregateGRPC(ctx, server.config, request)
	if err != nil {
		return nil, err
	}

	rsp, err := server.repo.ChangeStatusJobByAdmin(*changeStatusJobByAdminAggregate)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}
