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

type EditJobAggregate struct {
	JobID      uuid.UUID
	EmployerID string
	Job        map[string]interface{}
}

func EditJobAggregateGRPC(ctx context.Context, config utils.Config, request *pb.EditJobRequest) (*EditJobAggregate, error) {
	currentUser, err := middleware.AuthorizeUser(ctx, config, []string{utils.EmployerRole})
	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	violations := validateEditJobRequest(request)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	id, _ := uuid.Parse(request.GetId())
	enterpriseId, _ := uuid.Parse(request.GetEnterpriseId())

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
	job["SalaryLevelDisplay"] = request.GetSalaryLevelDisplay()
	job["ExactSalary"] = request.GetExactSalary()
	job["RangeSalary"] = request.GetRangeSalary()
	job["PaidPeriod"] = request.GetPaidPeriod()
	job["ExpiresAt"] = request.GetExpiresAt()

	job["EnterpriseID"] = enterpriseId
	job["EnterpriseName"] = request.GetEnterpriseName()
	job["EnterpriseAddress"] = request.GetEnterpriseAddress()

	if request.GetCrawl() { // check accesstoken.username of crawl service
		job["Crawl"] = request.GetCrawl()
		job["JobURL"] = request.GetJobUrl()
		job["JobSourceName"] = request.GetJobSourceName()
	}

	return &EditJobAggregate{
		JobID:      id,
		EmployerID: currentUser.Username,
		Job:        job,
	}, nil
}

func validateEditJobRequest(req *pb.EditJobRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if !utils.IsSupportedCurrency(req.GetCurrency()) {
		violations = append(violations, utils.FieldViolation("currency", fmt.Errorf("not support this currency")))
	}

	if _, err := uuid.Parse(req.GetId()); err != nil {
		violations = append(violations, utils.FieldViolation("id", fmt.Errorf("must be uuid")))
	}

	if req.GetEnterpriseId() != "" {
		if _, err := uuid.Parse(req.GetEnterpriseId()); err != nil {
			violations = append(violations, utils.FieldViolation("enterprise_id", fmt.Errorf("must be uuid")))
		}
	}
	return violations
}
