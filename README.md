# backend-job-service
This repo is job Management Service of jobstreet application backend.

## Deploy

1. ```docker build -t thanhquy1105/backend-jobstreet-job-service-prod:latest .```
2. ```docker push thanhquy1105/backend-jobstreet-job-service-prod```
3. ```docker pull thanhquy1105/backend-jobstreet-job-service-prod:latest```
4. ```docker run --name backend-jobstreet-job-service-prod --network jobstreet-network -p 8080:8080 -p 9090:9090 -e DB_SOURCE="postgresql://admin:admin@postgres:5432/job_service_jobstreet?sslmode=disable" -d thanhquy1105/backend-jobstreet-job-service-prod:latest```