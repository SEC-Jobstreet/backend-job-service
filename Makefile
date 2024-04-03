run_mongo:
	-docker network create job-network
	docker run --name mongo --network job-network -p 27017:27017 -d mongo 

start_mongo:
	docker start mongo

gql_generate:
	go run github.com/99designs/gqlgen generate

.PHONY: run_mongo gql_generate start_mongo