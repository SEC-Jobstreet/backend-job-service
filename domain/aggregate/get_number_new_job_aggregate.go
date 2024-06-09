package aggregate

import (
	"context"

	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

type GetNumberOfNewJobAggregate struct {
	Keyword string
	Address string
}

func GetNumberOfNewJobAggregateGRPC(ctx context.Context, request *pb.GetNumberOfNewJobRequest) *GetNumberOfNewJobAggregate {
	return &GetNumberOfNewJobAggregate{
		Keyword: request.GetKeyword(),
		Address: request.GetAddress(),
	}
}
