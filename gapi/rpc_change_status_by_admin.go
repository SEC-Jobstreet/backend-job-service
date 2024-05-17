package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/SEC-Jobstreet/backend-job-service/models"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) ChangeStatusJobByAdmin(ctx context.Context, request *pb.ChangeStatusJobByAdminRequest) (*pb.ChangeStatusJobByAdminResponse, error) {
	_, err := server.authorizeUser(ctx, server.config, []string{utils.Admin})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateChangeStatusJobByAdminRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	id, _ := uuid.Parse(request.GetId())

	job := map[string]interface{}{}

	job["ID"] = id
	job["Status"] = strings.ToUpper(request.GetStatus())

	err = server.store.Model(&models.Jobs{}).Where("id = ?", id).Updates(job).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to close job: %s", err)
	}

	fmt.Println(job)

	rsp := &pb.ChangeStatusJobByAdminResponse{
		Status: "OK",
	}

	return rsp, nil
}

func validateChangeStatusJobByAdminRequest(req *pb.ChangeStatusJobByAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if _, err := uuid.Parse(req.GetId()); err != nil {
		violations = append(violations, fieldViolation("id", fmt.Errorf("must be uuid")))
	}
	if !utils.IsSupportedStatus(strings.ToUpper(req.GetStatus())) {
		violations = append(violations, fieldViolation("status", fmt.Errorf("must be supported")))
	}
	return violations
}
