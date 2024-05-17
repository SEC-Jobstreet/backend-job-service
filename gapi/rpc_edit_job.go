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

func (server *Server) EditJob(ctx context.Context, request *pb.EditJobRequest) (*pb.EditJobResponse, error) {
	currentUser, err := server.authorizeUser(ctx, server.config, []string{utils.EmployerRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateEditJobRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, invalidArgumentError(violations)
	}

	enterpriseId, err := uuid.Parse(request.GetEnterpriseId())
	if err != nil {
		return nil, invalidArgumentError(violations)
	}

	job := map[string]interface{}{}

	job["ID"] = id
	job["EmployerID"] = currentUser.Username
	job["Status"] = "REVIEW"
	job["Title"] = request.GetTitle()
	job["Type"] = request.GetType()
	job["WorkWhenever"] = request.GetWorkWhenever()
	job["WorkShift"] = request.GetWorkShift()
	job["Description"] = request.GetDescription()
	job["Visa"] = request.GetVisa()
	job["Experience"] = request.GetExperience()
	job["StartDate"] = request.GetStartDate()
	job["Currency"] = request.GetCurrency()
	job["ExactSalary"] = request.GetExactSalary()
	job["RangeSalary"] = request.GetRangeSalary()
	job["ExpiresAt"] = request.GetExpiresAt()

	job["EnterpriseID"] = enterpriseId
	job["EnterpriseName"] = request.GetEnterpriseName()
	job["EnterpriseAddress"] = request.GetEnterpriseAddress()

	if request.GetCrawl() { // check accesstoken.username of crawl service
		job["Crawl"] = request.GetCrawl()
		job["JobURL"] = request.GetJobUrl()
		job["JobSourceName"] = request.GetJobSourceName()
	}

	err = server.store.Model(&models.Jobs{}).Clauses(clause.Returning{}).
		Where("id = ?", id).Where("employer_id = ?", currentUser.Username).Updates(job).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update job: %s", err)
	}

	rsp := &pb.EditJobResponse{
		Job: convertMapJob(job),
	}

	return rsp, nil
}

func validateEditJobRequest(req *pb.EditJobRequest) (violations []*errdetails.BadRequest_FieldViolation) {
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
