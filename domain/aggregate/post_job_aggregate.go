package aggregate

import (
	"context"
	"fmt"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate/entity"
	"github.com/SEC-Jobstreet/backend-job-service/domain/repository/model"
	"github.com/SEC-Jobstreet/backend-job-service/domain/utils"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type PostJobAggregate struct {
	Job        model.Jobs
	Enterprise entity.Enterprise
}

func PostJobAggregateGRPC(ctx context.Context, request *pb.PostJobRequest) (*PostJobAggregate, error) {
	violations := validatePostJobRequest(request)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, utils.InvalidArgumentError(violations)
	}

	enterpriseId, _ := uuid.NewRandom()

	job := model.Jobs{
		ID:           id,
		EmployerID:   request.GetEmployerId(),
		Status:       "REVIEW",
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
		ExpiresAt:    request.GetExpiresAt(),

		EnterpriseID:      enterpriseId,
		EnterpriseName:    request.GetEnterpriseName(),
		EnterpriseAddress: request.GetEnterpriseAddress(),
	}

	if request.GetCrawl() { // check accesstoken.username of crawl service
		job.Crawl = request.GetCrawl()
		job.JobURL = request.GetJobUrl()
		job.JobSourceName = request.GetJobSourceName()
	}

	enterprise := entity.Enterprise{
		ID:        enterpriseId,
		Name:      request.GetEnterpriseName(),
		Country:   request.GetEnterpriseCountry(),
		Address:   request.GetEnterpriseAddress(),
		Latitude:  request.GetEnterpriseLatitude(),
		Longitude: request.GetEnterpriseLongitude(),
		Field:     request.GetEnterpriseField(),
		Size:      request.GetEnterpriseSize(),
		Url:       request.GetEnterpriseUrl(),

		EmployerID:   request.GetEmployerId(),
		EmployerRole: request.GetEmployerRole(),
	}

	return &PostJobAggregate{
		Job:        job,
		Enterprise: enterprise,
	}, nil
}

func validatePostJobRequest(req *pb.PostJobRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if !utils.IsSupportedCurrency(req.GetCurrency()) {
		violations = append(violations, utils.FieldViolation("currency", fmt.Errorf("not support this currency")))
	}

	return violations
}
