package aggregate

import (
	"context"
	"fmt"

	"github.com/SEC-Jobstreet/backend-job-service/domain/utils"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type GetJobListAggregate struct {
	Keyword  string
	Address  string
	PageId   int32
	PageSize int32
}

func GetJobListAggregateGRPC(ctx context.Context, request *pb.JobListRequest) (*GetJobListAggregate, error) {

	violations := validateGetJobListRequest(request)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	return &GetJobListAggregate{
		Keyword:  request.GetKeyword(),
		Address:  request.GetAddress(),
		PageId:   request.GetPageId(),
		PageSize: request.GetPageSize(),
	}, nil
}

func validateGetJobListRequest(req *pb.JobListRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.GetPageId() < 1 {
		violations = append(violations, utils.FieldViolation("page_id", fmt.Errorf("page_id must be greater than 0")))
	}
	if req.GetPageSize() < 5 && req.GetPageSize() > 10 {
		violations = append(violations, utils.FieldViolation("page_size", fmt.Errorf("page_size must be from 5 to 10")))
	}
	return violations
}
