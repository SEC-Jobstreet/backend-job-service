package aggregate

import (
	"context"
	"fmt"
	"strings"

	"github.com/SEC-Jobstreet/backend-job-service/domain/middleware"
	"github.com/SEC-Jobstreet/backend-job-service/domain/utils"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type ChangeStatusJobByAdminAggregate struct {
	JobID uuid.UUID
	Job   map[string]interface{}
}

func ChangeStatusJobByAdminAggregateGRPC(ctx context.Context, config utils.Config, request *pb.ChangeStatusJobByAdminRequest) (*ChangeStatusJobByAdminAggregate, error) {
	_, err := middleware.AuthorizeUser(ctx, config, []string{utils.Admin})
	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	violations := validateChangeStatusJobByAdminRequest(request)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	id, _ := uuid.Parse(request.GetId())

	job := map[string]interface{}{}

	job["ID"] = id
	job["Status"] = strings.ToUpper(request.GetStatus())

	return &ChangeStatusJobByAdminAggregate{
		JobID: id,
		Job:   job,
	}, nil
}

func validateChangeStatusJobByAdminRequest(req *pb.ChangeStatusJobByAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if _, err := uuid.Parse(req.GetId()); err != nil {
		violations = append(violations, utils.FieldViolation("id", fmt.Errorf("must be uuid")))
	}
	if !utils.IsSupportedStatus(strings.ToUpper(req.GetStatus())) {
		violations = append(violations, utils.FieldViolation("status", fmt.Errorf("must be supported")))
	}
	return violations
}
