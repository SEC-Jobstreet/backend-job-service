run_mongo:
	-docker network create job-network
	docker run --name mongo --network job-network -p 27017:27017 -d mongo 

docker_compose_build:
	docker-compose build

build_app:
	docker build -t thanhquy1105/backend-jobstreet-job-service-prod:latest .

docker_push:
	docker push thanhquy1105/backend-jobstreet-job-service-prod

run_app:
	docker run --name backend-jobstreet-job-service-prod --network job-network -p 4001:4001 -e DB_URL="mongodb://mongo:27017" thanhquy1105/backend-jobstreet-job-service-prod:latest

start_app:
	docker start backend-jobstreet-job-service-prod

start_mongo:
	docker start mongo

gql_generate:
	go run github.com/99designs/gqlgen generate

.PHONY: run_mongo gql_generate start_mongo build_app docker_push run_app start_app