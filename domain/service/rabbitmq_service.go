package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SEC-Jobstreet/backend-job-service/domain/repository/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQService struct {
	serviceAddress string
	conn           *amqp.Connection
}

func NewRabbitMQService(address string, conn *amqp.Connection) *RabbitMQService {

	return &RabbitMQService{
		serviceAddress: address,
		conn:           conn,
	}
}

func (rs *RabbitMQService) PublishMessageToQueue(job model.Jobs) {
	ch, err := rs.conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"publishJob",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(q)

	data, err := json.Marshal(job)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()
	err = ch.PublishWithContext(ctx, "", "publishJob", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully Published Message to Queue")
}
