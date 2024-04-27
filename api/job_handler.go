package api

import (
	"net/http"

	"github.com/SEC-Jobstreet/backend-job-service/models"
	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type jobPostingRequest struct {
	Title        string `json:"title" binding:"required"`
	Type         string `json:"type"` // FULLTIME, PARTIME, SEASONAL
	WorkWhenever bool   `json:"work_whenever" `
	WorkShift    string `json:"work_shift"`
	Description  string `json:"description" binding:"required"`
	Visa         bool   `json:"visa"`
	Experience   uint   `json:"experience"` // 1: 0kn, 2: 1kn, 3: 23kn, 4: 4kn
	StartDate    int64  `json:"start_date" `
	Currency     string `json:"currency" binding:"currency"` // AUD, GBP, HKD, IDR, MYR, NZD, PHP, SGD, THB, VND
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

func (s *Server) PostJob(ctx *gin.Context) {

	var request jobPostingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	job := &models.Jobs{
		ID:           id,
		EmployerID:   currentUser.Username, // get accesstoken.username for employerID
		Status:       "POSTED",
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

	ctx.JSON(http.StatusOK, job)
}

type jobGettingRequest struct {
	JobID string `uri:"id" binding:"required"`
}

func (s *Server) GetJob(ctx *gin.Context) {
	var req jobGettingRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	id, err := uuid.Parse(req.JobID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	job := &models.Jobs{
		ID: id,
	}
	s.store.First(job)
	ctx.JSON(http.StatusOK, job)
}

type jobListRequest struct {
	Keyword  string `form:"keyword"`
	Address  string `form:"address"`
	PageID   int    `form:"page_id" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=10,max=20"`
}

type jobListResponse struct {
	Total int           `json:"total"`
	Data  []models.Jobs `json:"data"`
}

func (s *Server) JobList(ctx *gin.Context) {
	var req jobListRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	jobs := []models.Jobs{}

	keyword := "%" + req.Keyword + "%"
	address := "%" + req.Address + "%"

	tx := s.store.Where(s.store.Where(
		s.store.Where(
			s.store.Where("title LIKE ?", keyword),
		).Or(
			s.store.Where("description LIKE ?", keyword),
		)).Where(
		s.store.Where(
			s.store.Where("enterprise_address LIKE ?", address),
		).Or(
			s.store.Where("description LIKE ?", address),
		)).Where(
		s.store.Where("status = ?", "POSTED"),
	)).Find(&jobs)

	total := len(jobs)
	tx.Limit(req.PageSize).Offset((req.PageID - 1) * req.PageSize).Find(&jobs)

	res := jobListResponse{
		Total: total,
		Data:  jobs,
	}

	ctx.JSON(http.StatusOK, res)
}

type jobListByEmployerRequest struct {
	PageID   int `form:"page_id" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=10,max=20"`
}

func (s *Server) GetJobByEmployer(ctx *gin.Context) {
	var req jobListByEmployerRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	currentUser, err := utils.GetCurrentUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	jobs := []models.Jobs{}

	tx := s.store.Where("employer_id = ?", currentUser.Username).Find(&jobs)

	total := len(jobs)
	tx.Limit(req.PageSize).Offset((req.PageID - 1) * req.PageSize).Find(&jobs)

	res := jobListResponse{
		Total: total,
		Data:  jobs,
	}

	ctx.JSON(http.StatusOK, res)
}

func (s *Server) example(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}
