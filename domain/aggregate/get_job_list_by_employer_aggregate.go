package aggregate

import (
	"context"
	"fmt"

	"github.com/SEC-Jobstreet/backend-job-service/domain/middleware"
	"github.com/SEC-Jobstreet/backend-job-service/domain/utils"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type GetJobListByEmployerAggregate struct {
	EmployerID string
	PageId     int32
	PageSize   int32
}

func GetJobListByEmployerAggregateGRPC(ctx context.Context, config utils.Config, request *pb.JobListByEmployerRequest) (*GetJobListByEmployerAggregate, error) {
	currentUser, err := middleware.AuthorizeUser(ctx, config, []string{utils.EmployerRole, utils.Admin})
	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	violations := validateGetJobByEmployerRequest(request)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	return &GetJobListByEmployerAggregate{
		EmployerID: currentUser.Username,
		PageId:     request.GetPageId(),
		PageSize:   request.GetPageSize(),
	}, nil
}

func validateGetJobByEmployerRequest(req *pb.JobListByEmployerRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.GetPageId() < 1 {
		violations = append(violations, utils.FieldViolation("page_id", fmt.Errorf("page_id must be greater than 0")))
	}
	if req.GetPageSize() < 5 && req.GetPageSize() > 10 {
		violations = append(violations, utils.FieldViolation("page_size", fmt.Errorf("page_size must be from 5 to 10")))
	}
	return violations
}
