package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/SEC-Jobstreet/backend-job-service/domain/aggregate/entity"
)

type EmployerService struct {
	serviceAddress string
}

func NewEmployerService(serviceAddress string) *EmployerService {
	return &EmployerService{
		serviceAddress: serviceAddress,
	}
}

func (es *EmployerService) CreateEnterprise(employer entity.Enterprise) error {
	data, err := json.Marshal(employer)
	if err != nil {
		log.Fatal(err)
	}
	reader := bytes.NewReader(data)

	requestURL := fmt.Sprintf("%s/api/v1/create_enterprise", strings.TrimSpace(es.serviceAddress))
	req, err := http.NewRequest(http.MethodPost, requestURL, reader)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		return err
	}
	return err
}
