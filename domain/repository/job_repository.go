package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate"
	"github.com/SEC-Jobstreet/backend-job-service/domain/factory"
	"github.com/SEC-Jobstreet/backend-job-service/domain/repository/model"
	"github.com/SEC-Jobstreet/backend-job-service/domain/service"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JobRepository interface {
	Create(aggregate aggregate.PostJobAggregate) error
	GetJobByID(aggregate aggregate.GetJobByIDAggregate) (*pb.GetJobByIDResponse, error)
	GetJobList(aggregate aggregate.GetJobListAggregate) (*pb.JobListResponse, error)
	GetNumberOfJob() *pb.GetNumberOfJobResponse
	GetNumberOfNewJob(aggregate aggregate.GetNumberOfNewJobAggregate) *pb.GetNumberOfNewJobResponse
	GetJobListByEmployer(aggregate aggregate.GetJobListByEmployerAggregate) (*pb.JobListResponse, error)
	GetJobListByAdmin(aggregate aggregate.GetJobListByAdminAggregate) (*pb.JobListResponse, error)
	EditJob(aggregate aggregate.EditJobAggregate) (*pb.EditJobResponse, error)
	CloseJob(aggregate aggregate.CloseJobAggregate) (*pb.CloseJobResponse, error)
	ChangeStatusJobByAdmin(aggregate aggregate.ChangeStatusJobByAdminAggregate) (*pb.ChangeStatusJobByAdminResponse, error)
}

type jobRepo struct {
	db              *gorm.DB
	rdb             *RedisJobRepository
	employerService *service.EmployerService
	rabbitmq        *service.RabbitMQService
}

func NewJobRepository(db *gorm.DB, rdb *RedisJobRepository, es *service.EmployerService, rabbitmq *service.RabbitMQService) JobRepository {
	return &jobRepo{
		db:              db,
		rdb:             rdb,
		employerService: es,
		rabbitmq:        rabbitmq,
	}
}

func (jr *jobRepo) Create(aggregate aggregate.PostJobAggregate) error {
	tx := jr.db.Begin()
	tx.Create(&aggregate.Job)
	if tx.Error != nil {
		tx.Rollback()
		if tx.Error == gorm.ErrDuplicatedKey || tx.Error == gorm.ErrForeignKeyViolated {
			return status.Errorf(codes.AlreadyExists, "%s", tx.Error)
		}
		return status.Errorf(codes.Internal, "failed to create account: %s", tx.Error)
	}

	err := jr.employerService.CreateEnterprise(aggregate.Enterprise)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (jr *jobRepo) GetJobList(aggregate aggregate.GetJobListAggregate) (*pb.JobListResponse, error) {
	// Get from redis
	key := "jobs?keyword=" + aggregate.Keyword + "&address=" + aggregate.Address +
		"&pageid=" + fmt.Sprintf("%v", aggregate.PageId) + "&pagesize=" + fmt.Sprintf("%v", aggregate.PageSize)
	if jobsListRedis, err := jr.rdb.GetJobList(context.Background(), key); err == nil && len(jobsListRedis.Jobs) != 0 {
		jobsListRedis.PageId = aggregate.PageId
		jobsListRedis.PageSize = aggregate.PageSize
		return jobsListRedis, err
	}

	// Get from db
	jobs := []model.Jobs{}

	keyword := "%" + aggregate.Keyword + "%"
	address := "%" + aggregate.Address + "%"

	var total int64

	tx := jr.db.Model(&model.Jobs{}).Where(
		jr.db.Where(
			jr.db.Where(
				jr.db.Where("title LIKE ?", keyword),
			).Or(
				jr.db.Where("description LIKE ?", keyword),
			)).Where(
			jr.db.Where(
				jr.db.Where("enterprise_address LIKE ?", address),
			).Or(
				jr.db.Where("description LIKE ?", address),
			)).
			Where("status = ?", "POSTED"))
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get jobs:%s", tx.Error)
	}

	tx.Count(&total)

	err := tx.Limit(int(aggregate.PageSize)).Offset(int(aggregate.PageId-1) * int(aggregate.PageSize)).Find(&jobs).Order("created_at DESC").Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get jobs:%s", err)
	}

	res := pb.JobListResponse{
		Total:    total,
		PageId:   aggregate.PageId,
		PageSize: aggregate.PageSize,
		Jobs:     factory.ConvertJobList(jobs),
	}

	// Save to redis
	jr.rdb.SaveJobList(context.Background(), key, total, jobs)

	return &res, err
}

func (jr *jobRepo) GetJobListByEmployer(aggregate aggregate.GetJobListByEmployerAggregate) (*pb.JobListResponse, error) {
	jobs := []model.Jobs{}
	var total int64

	tx := jr.db.Model(&model.Jobs{}).Where("employer_id = ?", aggregate.EmployerID)
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get jobs:%s", tx.Error)
	}
	tx.Count(&total)

	err := tx.Limit(int(aggregate.PageSize)).Offset(int(aggregate.PageId-1) * int(aggregate.PageSize)).Find(&jobs).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get jobs:%s", err)
	}

	return &pb.JobListResponse{
		Total:    total,
		PageId:   aggregate.PageId,
		PageSize: aggregate.PageSize,
		Jobs:     factory.ConvertJobList(jobs),
	}, nil
}

func (jr *jobRepo) GetJobListByAdmin(aggregate aggregate.GetJobListByAdminAggregate) (*pb.JobListResponse, error) {
	jobs := []model.Jobs{}
	var total int64

	tx := jr.db.Model(&model.Jobs{}).Order("status DESC").Order("updated_at DESC")
	if aggregate.Status != "" {
		tx = tx.Where("status = ?", strings.ToUpper(aggregate.Status))
	}
	if tx.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to get jobs:%s", tx.Error)
	}

	tx.Count(&total)

	err := tx.Limit(int(aggregate.PageSize)).Offset(int(aggregate.PageId-1) * int(aggregate.PageSize)).Find(&jobs).Error

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get jobs:%s", err)
	}

	return &pb.JobListResponse{
		Total:    total,
		PageId:   aggregate.PageId,
		PageSize: aggregate.PageSize,
		Jobs:     factory.ConvertJobList(jobs),
	}, nil
}

func (jr *jobRepo) GetJobByID(aggregate aggregate.GetJobByIDAggregate) (*pb.GetJobByIDResponse, error) {
	// Get from redis
	jobRdis, err := jr.rdb.GetJob(context.Background(), "job/"+aggregate.ID.String())
	if err == nil {
		return &pb.GetJobByIDResponse{
			Job: factory.ConvertJob(*jobRdis),
		}, nil
	}

	// Get from db
	job := &model.Jobs{
		ID: aggregate.ID,
	}
	err = jr.db.First(job).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get job:%s", err)
	}

	// save to redis
	jr.rdb.SaveJob(context.Background(), "job/"+aggregate.ID.String(), *job)

	return &pb.GetJobByIDResponse{
		Job: factory.ConvertJob(*job),
	}, nil
}

func (jr *jobRepo) EditJob(aggregate aggregate.EditJobAggregate) (*pb.EditJobResponse, error) {
	err := jr.db.Model(&model.Jobs{}).Clauses(clause.Returning{}).
		Where("id = ?", aggregate.JobID).Where("employer_id = ?", aggregate.EmployerID).Updates(aggregate.Job).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update job: %s", err)
	}

	return &pb.EditJobResponse{
		Job: factory.ConvertMapJob(aggregate.Job),
	}, nil
}

func (jr *jobRepo) CloseJob(aggregate aggregate.CloseJobAggregate) (*pb.CloseJobResponse, error) {

	err := jr.db.Model(&model.Jobs{}).Clauses(clause.Returning{}).
		Where("id = ?", aggregate.JobID).Where("employer_id = ?", aggregate.EmployerID).Updates(aggregate.Job).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to close job: %s", err)
	}

	return &pb.CloseJobResponse{
		Status: "OK",
	}, nil
}

func (jr *jobRepo) ChangeStatusJobByAdmin(aggregate aggregate.ChangeStatusJobByAdminAggregate) (*pb.ChangeStatusJobByAdminResponse, error) {

	err := jr.db.Model(&model.Jobs{}).Where("id = ?", aggregate.JobID).Updates(aggregate.Job).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to close job: %s", err)
	}

	go func() {
		if aggregate.Job["Status"] == "POSTED" {
			job := &model.Jobs{
				ID: aggregate.JobID,
			}
			err = jr.db.First(job).Error
			if err != nil {
				fmt.Println(status.Errorf(codes.Internal, "failed to get job:%s", err))
				return
			}
			jr.rabbitmq.PublishMessageToQueue(*job)
		}
	}()

	return &pb.ChangeStatusJobByAdminResponse{
		Status: "OK",
	}, nil
}

func (jr *jobRepo) GetNumberOfJob() *pb.GetNumberOfJobResponse {
	if total, err := jr.rdb.GetNumber(context.Background(), "numberofjob"); err == nil && total != 0 {
		return &pb.GetNumberOfJobResponse{
			Total: total,
		}
	}

	var total int64
	jr.db.Model(&model.Jobs{}).Where("status = ?", "POSTED").Count(&total)

	jr.rdb.SaveNumber(context.Background(), "numberofjob", total)

	return &pb.GetNumberOfJobResponse{
		Total: total,
	}
}

func (jr *jobRepo) GetNumberOfNewJob(aggregate aggregate.GetNumberOfNewJobAggregate) *pb.GetNumberOfNewJobResponse {
	key := "numberofnewjob?" + "keyword=" + aggregate.Keyword + "&address=" + aggregate.Address
	if total, err := jr.rdb.GetNumber(context.Background(), key); err == nil && total != 0 {
		return &pb.GetNumberOfNewJobResponse{
			Total: total,
		}
	}

	var total int64
	t := time.Now().Add(-24 * time.Hour).Unix()

	jr.db.Model(&model.Jobs{}).Where("status = ?", "POSTED").Where("updated_at > ?", t).Count(&total)

	jr.rdb.SaveNumber(context.Background(), key, total)

	return &pb.GetNumberOfNewJobResponse{
		Total: total,
	}
}
