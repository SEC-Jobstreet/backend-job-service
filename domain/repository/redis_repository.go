package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/SEC-Jobstreet/backend-job-service/domain/factory"
	"github.com/SEC-Jobstreet/backend-job-service/domain/repository/model"
	"github.com/SEC-Jobstreet/backend-job-service/pb"
	"github.com/redis/go-redis/v9"
)

type RedisJobRepository struct {
	rdb *redis.Client
}

func NewRedisJobRepository(rdb *redis.Client) *RedisJobRepository {
	return &RedisJobRepository{
		rdb: rdb,
	}
}

type jobList struct {
	Total int64
	Jobs  []model.Jobs
}

func (rr *RedisJobRepository) SaveJobList(ctx context.Context, key string, total int64, value []model.Jobs) error {
	data, err := json.Marshal(jobList{
		Total: total,
		Jobs:  value,
	})
	if err != nil {
		return err
	}

	err = rr.rdb.Set(ctx, key, data, time.Duration(24)*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rr *RedisJobRepository) GetJobList(ctx context.Context, key string) (*pb.JobListResponse, error) {
	value, err := rr.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var retrievedJob jobList
	err = json.Unmarshal([]byte(value), &retrievedJob)
	if err != nil {
		return nil, err
	}

	return &pb.JobListResponse{
		Total: retrievedJob.Total,
		Jobs:  factory.ConvertJobList(retrievedJob.Jobs),
	}, nil
}

func (rr *RedisJobRepository) SaveJob(ctx context.Context, key string, value model.Jobs) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = rr.rdb.Set(ctx, key, data, time.Duration(24)*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rr *RedisJobRepository) GetJob(ctx context.Context, key string) (*model.Jobs, error) {
	value, err := rr.rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var retrievedJob model.Jobs
	err = json.Unmarshal([]byte(value), &retrievedJob)
	if err != nil {
		return nil, err
	}

	return &retrievedJob, nil
}

func (rr *RedisJobRepository) SaveNumber(ctx context.Context, key string, value int64) error {
	err := rr.rdb.Set(ctx, key, value, time.Duration(24)*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rr *RedisJobRepository) GetNumber(ctx context.Context, key string) (int64, error) {
	res, err := rr.rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(res)
	if err != nil {
		return 0, err
	}

	return int64(value), nil
}
