package aggregate

import (
	"context"
	"fmt"

	"github.com/SEC-Jobstreet/backend-job-service/domain/utils"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type GetJobByIDAggregate struct {
	ID uuid.UUID
}

func GetJobByIDAggregateGRPC(ctx context.Context, request *pb.GetJobByIDRequest) (*GetJobByIDAggregate, error) {
	violations := validateGetJobByIDRequest(request)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	id, _ := uuid.Parse(request.GetId())

	return &GetJobByIDAggregate{
		ID: id,
	}, nil
}

func validateGetJobByIDRequest(req *pb.GetJobByIDRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if _, err := uuid.Parse(req.GetId()); err != nil {
		violations = append(violations, utils.FieldViolation("id", fmt.Errorf("must be uuid")))
	}
	return violations
}
