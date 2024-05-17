package gapi

import (
	"context"
	"fmt"

	"github.com/SEC-Jobstreet/backend-job-service/models"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm/clause"
)

func (server *Server) CloseJob(ctx context.Context, request *pb.CloseJobRequest) (*pb.CloseJobResponse, error) {
	currentUser, err := server.authorizeUser(ctx, server.config, []string{utils.EmployerRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateCloseJobRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	id, _ := uuid.Parse(request.GetId())

	job := map[string]interface{}{}

	job["ID"] = id
	job["EmployerID"] = currentUser.Username
	job["Status"] = "CLOSED"

	err = server.store.Model(&models.Jobs{}).Clauses(clause.Returning{}).
		Where("id = ?", id).Where("employer_id = ?", currentUser.Username).Updates(job).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to close job: %s", err)
	}

	fmt.Println(job)

	rsp := &pb.CloseJobResponse{
		Status: "OK",
	}

	return rsp, nil
}

func validateCloseJobRequest(req *pb.CloseJobRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if _, err := uuid.Parse(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", fmt.Errorf("must be uuid")))
	}
	return violations
}
