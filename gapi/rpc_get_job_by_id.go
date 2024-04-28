package gapi

import (
	"context"
	"fmt"

	"github.com/SEC-Jobstreet/backend-job-service/models"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetJobByID(ctx context.Context, request *pb.GetJobByIDRequest) (*pb.GetJobByIDResponse, error) {
	violations := validateGetJobByIDRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	id, _ := uuid.Parse(request.GetId())

	job := &models.Jobs{
		ID: id,
	}
	err := server.store.First(job).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get job:%s", err)
	}

	rsp := &pb.GetJobByIDResponse{
		Job: convertJob(*job),
	}

	return rsp, nil
}

func validateGetJobByIDRequest(req *pb.GetJobByIDRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if _, err := uuid.Parse(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", fmt.Errorf("must be uuid")))
	}
	return violations
}
