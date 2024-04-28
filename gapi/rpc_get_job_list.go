package gapi

import (
	"context"
	"fmt"

	"github.com/SEC-Jobstreet/backend-job-service/models"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetJobList(ctx context.Context, request *pb.JobListRequest) (*pb.JobListResponse, error) {

	violations := validateGetJobListRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	jobs := []models.Jobs{}

	keyword := "%" + request.GetKeyword() + "%"
	address := "%" + request.GetAddress() + "%"

	var total int64

	tx := server.store.Model(&models.Jobs{}).Where(
		server.store.Where(
			server.store.Where(
				server.store.Where("title LIKE ?", keyword),
			).Or(
				server.store.Where("description LIKE ?", keyword),
			)).Where(
			server.store.Where(
				server.store.Where("enterprise_address LIKE ?", address),
			).Or(
				server.store.Where("description LIKE ?", address),
			)).
			Where("status = ?", "POSTED"))
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get jobs:%s", tx.Error)
	}

	tx.Count(&total)

	err := tx.Limit(int(request.GetPageSize())).Offset(int(request.GetPageId()-1) * int(request.PageSize)).Find(&jobs).Error
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

func validateGetJobListRequest(req *pb.JobListRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.GetPageId() < 1 {
		violations = append(violations, fieldViolation("page_id", fmt.Errorf("page_id must be greater than 0")))
	}
	if req.GetPageSize() < 5 && req.GetPageSize() > 10 {
		violations = append(violations, fieldViolation("page_size", fmt.Errorf("page_size must be from 5 to 10")))
	}
	return violations
}
