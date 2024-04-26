package api

import (
	"fmt"
	"net/http"

	"github.com/SEC-Jobstreet/backend-job-service/models"
	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type jobPostingRequest struct {
	// EmployerID   uint      `json:"employer_id"`
	Title        string `json:"title" binding:"required"`
	Type         string `json:"type"` // FULLTIME, PARTIME, SEASONAL
	WorkWhenever bool   `json:"work_whenever" `
	WorkShift    string `json:"work_shift"`
	Description  string `json:"description" binding:"required"`
	Visa         bool   `json:"visa"`
	Experience   uint   `json:"experience"` // 1: 0kn, 2: 1kn, 3: 23kn, 4: 4kn
	StartDate    int64  `json:"start_date" `
	Currency     string `json:"currency"` // AUD, GBP, HKD, IDR, MYR, NZD, PHP, SGD, THB, VND
	ExactSalary  uint   `json:"exact_salary"`
	RangeSalary  string `json:"range_salary"`
	ExpireAt     int64  `json:"expire_at" `

	EnterpriseID      uuid.UUID `json:"enterprise_id"`
	EnterpriseName    string    `json:"enterprise_name"`
	EnterpriseAddress string    `json:"enterprise_address"`

	Crawl         bool   `json:"crawl"`
	JobURL        string `json:"job_url"`
	JobSourceName string `json:"job_source_name"`
}

type jobPostingResponse struct {
	Status string `json:"status"`
	Job    jobPostingRequest
}

func (s *Server) PostJob(ctx *gin.Context) {
	var request jobPostingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	fmt.Println(request)

	ctx.JSON(http.StatusOK, request)

	id, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	job := &models.Jobs{
		ID:           id,
		EmployerID:   "test", // get accesstoken.username for employerID
		Status:       "REVIEW",
		Title:        request.Title,
		Type:         request.Type,
		WorkWhenever: request.WorkWhenever,
		WorkShift:    request.WorkShift,
		Description:  request.Description,
		Visa:         request.Visa,
		Experience:   request.Experience,
		StartDate:    request.StartDate,
		Currency:     request.Currency,
		ExactSalary:  request.ExactSalary,
		RangeSalary:  request.RangeSalary,
		ExpireAt:     request.ExpireAt,

		EnterpriseID:      request.EnterpriseID,
		EnterpriseName:    request.EnterpriseName,
		EnterpriseAddress: request.EnterpriseAddress,
	}

	if request.Crawl { // check accesstoken.username of crawl service
		job.Crawl = request.Crawl
		job.JobURL = request.JobURL
		job.JobSourceName = request.JobSourceName
	}

	s.store.Create(job)
}

func (s *Server) example(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, "OK")
}
