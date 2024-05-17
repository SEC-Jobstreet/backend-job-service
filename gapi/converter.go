package gapi

import (
	"fmt"

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
		ExpiresAt:    job.ExpiresAt,

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

func convertMapJob(job map[string]interface{}) (res *pb.Job) {
	res = &pb.Job{
		Id:          fmt.Sprintf("%v", job["ID"]),
		EmployerId:  fmt.Sprintf("%v", job["EmployerID"]),
		Status:      fmt.Sprintf("%v", job["Status"]),
		Title:       fmt.Sprintf("%v", job["Title"]),
		Type:        fmt.Sprintf("%v", job["Type"]),
		WorkShift:   fmt.Sprintf("%v", job["WorkShift"]),
		Description: fmt.Sprintf("%v", job["Description"]),

		Currency:    fmt.Sprintf("%v", job["Currency"]),
		RangeSalary: fmt.Sprintf("%v", job["RangeSalary"]),

		EnterpriseId:      fmt.Sprintf("%v", job["EnterpriseID"]),
		EnterpriseName:    fmt.Sprintf("%v", job["EnterpriseName"]),
		EnterpriseAddress: fmt.Sprintf("%v", job["EnterpriseAddress"]),

		JobUrl:        fmt.Sprintf("%v", job["JobURL"]),
		JobSourceName: fmt.Sprintf("%v", job["JobSourceName"]),
	}

	if job["Visa"] != nil {
		res.Visa = job["Visa"].(bool)
	}
	if job["WorkWhenever"] != nil {
		res.WorkWhenever = job["WorkWhenever"].(bool)
	}
	if job["Experience"] != nil {
		res.Experience = job["Experience"].(uint32)
	}
	if job["StartDate"] != nil {
		res.StartDate = job["StartDate"].(int64)
	}
	if job["CreatedAt"] != nil {
		res.CreatedAt = job["CreatedAt"].(int64)
	}
	if job["ExactSalary"] != nil {
		res.ExactSalary = job["ExactSalary"].(uint32)
	}
	if job["ExpiresAt"] != nil {
		res.ExpiresAt = job["ExpiresAt"].(int64)
	}
	if job["UpdatedAt"] != nil {
		res.UpdatedAt = job["UpdatedAt"].(int64)
	}
	if job["Crawl"] != nil {
		res.Crawl = job["Crawl"].(bool)
	}

	return res
}
