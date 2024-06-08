package aggregate

import (
	"context"
	"fmt"

	"github.com/SEC-Jobstreet/backend-job-service/domain/middleware"
	"github.com/SEC-Jobstreet/backend-job-service/domain/utils"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type GetJobListByAdminAggregate struct {
	PageId   int32
	PageSize int32
	Status   string
}

func GetJobListByAdminAggregateGRPC(ctx context.Context, config utils.Config, request *pb.JobListByAdminRequest) (*GetJobListByAdminAggregate, error) {
	_, err := middleware.AuthorizeUser(ctx, config, []string{utils.Admin})
	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	violations := validateGetJobByAdminRequest(request)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	return &GetJobListByAdminAggregate{
		PageId:   request.GetPageId(),
		PageSize: request.GetPageSize(),
		Status:   request.GetStatus(),
	}, nil
}

func validateGetJobByAdminRequest(req *pb.JobListByAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.GetPageId() < 1 {
		violations = append(violations, utils.FieldViolation("page_id", fmt.Errorf("page_id must be greater than 0")))
	}
	if req.GetPageSize() < 5 && req.GetPageSize() > 10 {
		violations = append(violations, utils.FieldViolation("page_size", fmt.Errorf("page_size must be from 5 to 10")))
	}
	if req.GetStatus() != "" {
		if !utils.IsSupportedStatus(req.GetStatus()) {
			violations = append(violations, utils.FieldViolation("status", fmt.Errorf("must be supported")))
		}
	}
	return violations
}
