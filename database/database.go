package database

import (
	"context"
	"log"
	"time"

	"github.com/SEC-Jobstreet/backend-job-service/graph/model"
	"github.com/SEC-Jobstreet/backend-job-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client *mongo.Client
}

func Connect(config *utils.Config) *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.DB_URL))
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		client: client,
	}
}

func (db *DB) GetJob(id string) *model.JobListing {
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var jobListing model.JobListing
	err := jobCollec.FindOne(ctx, filter).Decode(&jobListing)
	if err != nil {
		log.Fatal(err)
	}
	return &jobListing
}

func (db *DB) GetJobs() []*model.JobListing {
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var jobListings []*model.JobListing
	cursor, err := jobCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &jobListings); err != nil {
		panic(err)
	}

	return jobListings
}

func (db *DB) CreateJobListing(jobInfo model.CreateJobListingInput) *model.JobListing {
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserg, err := jobCollec.InsertOne(ctx, bson.M{"title": jobInfo.Title, "description": jobInfo.Description, "status": jobInfo.Status, "company": jobInfo.Company})

	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnJobListing := model.JobListing{ID: insertedID, Title: jobInfo.Title, Company: jobInfo.Company, Description: jobInfo.Description, Status: jobInfo.Status}
	return &returnJobListing
}

func (db *DB) UpdateJobListing(jobId string, jobInfo model.UpdateJobListingInput) *model.JobListing {
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateJobInfo := bson.M{}

	if jobInfo.Title != nil {
		updateJobInfo["title"] = jobInfo.Title
	}
	if jobInfo.Description != nil {
		updateJobInfo["description"] = jobInfo.Description
	}
	if jobInfo.Status != nil {
		updateJobInfo["status"] = jobInfo.Status
	}

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateJobInfo}

	results := jobCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var jobListing model.JobListing
	results.Decode(&jobListing)
	// if err := results.Decode(&jobListing); err != nil {
	// 	log.Fatal(err)
	// }

	return &jobListing
}

func (db *DB) DeleteJobListing(jobId string) *model.DeleteJobResponse {
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	_, err := jobCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	return &model.DeleteJobResponse{DeletedJobID: jobId}
}

func (db *DB) CreateSavedJob(jobID string, candidateID string) *model.SavedJob {
	jobCollec := db.client.Database("graphql-job-board").Collection("savedJobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	inserg, err := jobCollec.InsertOne(ctx, bson.M{"jobID": jobID, "candidateID": candidateID})
	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnSavedJob := model.SavedJob{ID: insertedID, JobID: jobID, CandidateID: candidateID}
	return &returnSavedJob
}

func (db *DB) GetSavedJobs(candidateID string) []*model.JobListing {
	savedjobCollec := db.client.Database("graphql-job-board").Collection("savedJobs")
	jobCollec := db.client.Database("graphql-job-board").Collection("jobs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// _candidateID, _ := primitive.ObjectIDFromHex(candidateID)
	filter := bson.M{"candidateID": candidateID}
	var savedJobs []*model.SavedJob
	cursor, err := savedjobCollec.Find(ctx, filter)

	if err = cursor.All(context.TODO(), &savedJobs); err != nil {
		panic(err)
	}

	var jobListings []*model.JobListing

	for i := 0; i < len(savedJobs); i++ {
		_id, _ := primitive.ObjectIDFromHex(savedJobs[i].JobID)
		filter := bson.M{"_id": _id}
		var jobListing model.JobListing
		err := jobCollec.FindOne(ctx, filter).Decode(&jobListing)
		if err != nil {
			log.Fatal("eeeeeeee")
		}
		jobListings = append(jobListings, &jobListing)
	}
	return jobListings
}
