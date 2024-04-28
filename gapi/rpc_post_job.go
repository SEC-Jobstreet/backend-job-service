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
	"gorm.io/gorm"
)

func (server *Server) PostJob(ctx context.Context, request *pb.PostJobRequest) (*pb.PostJobResponse, error) {
	currentUser, err := server.authorizeUser(ctx, server.config, []string{utils.EmployerRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validatePostJobRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, invalidArgumentError(violations)
	}

	enterpriseId, _ := uuid.Parse(request.GetEnterpriseId())

	job := &models.Jobs{
		ID:           id,
		EmployerID:   currentUser.Username, // get accesstoken.username for employerID
		Status:       "POSTED",
		Title:        request.GetTitle(),
		Type:         request.GetType(),
		WorkWhenever: request.GetWorkWhenever(),
		WorkShift:    request.GetWorkShift(),
		Description:  request.GetDescription(),
		Visa:         request.GetVisa(),
		Experience:   request.GetExperience(),
		StartDate:    request.GetStartDate(),
		Currency:     request.GetCurrency(),
		ExactSalary:  request.GetExactSalary(),
		RangeSalary:  request.GetRangeSalary(),
		ExpireAt:     request.GetExpiresAt(),

		EnterpriseID:      enterpriseId,
		EnterpriseName:    request.GetEnterpriseName(),
		EnterpriseAddress: request.GetEnterpriseAddress(),
	}

	if request.GetCrawl() { // check accesstoken.username of crawl service
		job.Crawl = request.GetCrawl()
		job.JobURL = request.GetJobUrl()
		job.JobSourceName = request.GetJobSourceName()
	}

	err = server.store.Create(job).Error
	if err != nil {
		if err == gorm.ErrDuplicatedKey || err == gorm.ErrForeignKeyViolated {
			return nil, status.Errorf(codes.AlreadyExists, "%s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create account: %s", err)
	}

	rsp := &pb.PostJobResponse{
		Job: convertJob(*job),
	}

	return rsp, nil
}

func validatePostJobRequest(req *pb.PostJobRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if !utils.IsSupportedCurrency(req.GetCurrency()) {
		violations = append(violations, fieldViolation("currency", fmt.Errorf("not support this currency")))
	}

	if req.GetEnterpriseId() != "" {
		if _, err := uuid.Parse(req.GetEnterpriseId()); err != nil {
			violations = append(violations, fieldViolation("enterprise_id", fmt.Errorf("must be uuid")))
		}
	}
	return violations
}
