package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/SEC-Jobstreet/backend-job-service/models"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetJobListByAdmin(ctx context.Context, request *pb.JobListByAdminRequest) (*pb.JobListResponse, error) {
	_, err := server.authorizeUser(ctx, server.config, []string{utils.Admin})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateGetJobByAdminRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	jobs := []models.Jobs{}
	var total int64

	tx := server.store.Model(&models.Jobs{}).Order("status DESC").Order("updated_at DESC")
	fmt.Println(request.GetStatus())
	if request.GetStatus() != "" {
		tx = tx.Where("status = ?", strings.ToUpper(request.GetStatus()))
	}
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get jobs:%s", tx.Error)
	}

	tx.Count(&total)

	err = tx.Limit(int(request.GetPageSize())).Offset(int(request.GetPageId()-1) * int(request.PageSize)).Find(&jobs).Error

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get jobs:%s", err)
	}

	rsp := &pb.JobListResponse{
		Total:    total,
		PageId:   request.GetPageId(),
		PageSize: request.GetPageSize(),
		Jobs:     convertJobList(jobs),
	}

	return rsp, nil
}

func validateGetJobByAdminRequest(req *pb.JobListByAdminRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.GetPageId() < 1 {
		violations = append(violations, fieldViolation("page_id", fmt.Errorf("page_id must be greater than 0")))
	}
	if req.GetPageSize() < 5 && req.GetPageSize() > 10 {
		violations = append(violations, fieldViolation("page_size", fmt.Errorf("page_size must be from 5 to 10")))
	}
	if req.GetStatus() != "" {
		if !utils.IsSupportedStatus(req.GetStatus()) {
			violations = append(violations, fieldViolation("status", fmt.Errorf("must be supported")))
		}
	}
	return violations
}
