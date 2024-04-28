package gapi

import (
	"github.com/SEC-Jobstreet/backend-job-service/models"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
)

func convertJob(job models.Jobs) *pb.Job {
	return &pb.Job{
		Id:           job.ID.String(),
		EmployerId:   job.EmployerID,
		Status:       job.Status,
		Title:        job.Title,
		Type:         job.Type,
		WorkWhenever: job.WorkWhenever,
		WorkShift:    job.WorkShift,
		Description:  job.Description,
		Visa:         job.Visa,
		Experience:   job.Experience,
		StartDate:    job.StartDate,
		Currency:     job.Currency,
		ExactSalary:  job.ExactSalary,
		RangeSalary:  job.RangeSalary,
		ExpiresAt:    job.ExpireAt,

		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,

		EnterpriseId:      job.EnterpriseID.String(),
		EnterpriseName:    job.EnterpriseName,
		EnterpriseAddress: job.EnterpriseAddress,

		Crawl:         job.Crawl,
		JobUrl:        job.JobURL,
		JobSourceName: job.JobSourceName,
	}
}

func convertJobList(jobs []models.Jobs) (res []*pb.Job) {
	for _, value := range jobs {
		res = append(res, convertJob(value))
	}
	return res
}
