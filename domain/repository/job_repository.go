package repository

import (
	"strings"

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
	GetJobList(aggregate aggregate.GetJobListAggregate) (*pb.JobListResponse, error)
	GetJobListByEmployer(aggregate aggregate.GetJobListByEmployerAggregate) (*pb.JobListResponse, error)
	GetJobListByAdmin(aggregate aggregate.GetJobListByAdminAggregate) (*pb.JobListResponse, error)
	GetJobByID(aggregate aggregate.GetJobByIDAggregate) (*pb.GetJobByIDResponse, error)
	EditJob(aggregate aggregate.EditJobAggregate) (*pb.EditJobResponse, error)
	CloseJob(aggregate aggregate.CloseJobAggregate) (*pb.CloseJobResponse, error)
	ChangeStatusJobByAdmin(aggregate aggregate.ChangeStatusJobByAdminAggregate) (*pb.ChangeStatusJobByAdminResponse, error)
}

type jobRepo struct {
	db              *gorm.DB
	employerService *service.EmployerService
}

func NewJobRepository(db *gorm.DB, es *service.EmployerService) JobRepository {
	return &jobRepo{
		db:              db,
		employerService: es,
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

	err := tx.Limit(int(aggregate.PageSize)).Offset(int(aggregate.PageId-1) * int(aggregate.PageSize)).Find(&jobs).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get jobs:%s", err)
	}

	return &pb.JobListResponse{
		Total:    total,
		PageId:   aggregate.PageId,
		PageSize: aggregate.PageSize,
		Jobs:     factory.ConvertJobList(jobs),
	}, err
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
	job := &model.Jobs{
		ID: aggregate.ID,
	}
	err := jr.db.First(job).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get job:%s", err)
	}

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

	return &pb.ChangeStatusJobByAdminResponse{
		Status: "OK",
	}, nil
}
