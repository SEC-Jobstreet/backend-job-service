package aggregate

import (
	"context"
	"fmt"

	"github.com/SEC-Jobstreet/backend-job-service/domain/middleware"
	"github.com/SEC-Jobstreet/backend-job-service/domain/utils"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type CloseJobAggregate struct {
	JobID      uuid.UUID
	EmployerID string
	Job        map[string]interface{}
}

func CloseJobAggregateGRPC(ctx context.Context, config utils.Config, request *pb.CloseJobRequest) (*CloseJobAggregate, error) {
	currentUser, err := middleware.AuthorizeUser(ctx, config, []string{utils.EmployerRole})
	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	violations := validateCloseJobRequest(request)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	id, _ := uuid.Parse(request.GetId())

	job := map[string]interface{}{}

	job["ID"] = id
	job["EmployerID"] = currentUser.Username
	job["Status"] = "CLOSED"

	return &CloseJobAggregate{
		JobID:      id,
		EmployerID: currentUser.Username,
		Job:        job,
	}, nil
}

func validateCloseJobRequest(req *pb.CloseJobRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if _, err := uuid.Parse(req.GetId()); err != nil {
		violations = append(violations, utils.FieldViolation("id", fmt.Errorf("must be uuid")))
	}
	return violations
}
