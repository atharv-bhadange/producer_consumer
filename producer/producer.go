package producer

import (
	"context"
	"fmt"

	"github.com/atharv-bhadange/producer_consumer/configs"
	"github.com/wagslane/go-rabbitmq"
)

type Qinstance struct {
	Conn *rabbitmq.Conn
}

var Rabbit Qinstance

func ConnectQueue() error {

	connString := configs.GetRabbitMQUrl()

	conn, err := rabbitmq.NewConn(connString)

	Rabbit = Qinstance{
		Conn: conn,
	}

	if err != nil {
		return err
	}

	return nil

}

func PublishMessage(product_id int) error {
	publisher, err := rabbitmq.NewPublisher(
		Rabbit.Conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)

	if err != nil {
		return err
	}

	defer publisher.Close()

	data := fmt.Sprintf("%v", product_id)

	err = publisher.PublishWithContext(
		context.Background(),
		[]byte(data),
		[]string{"secret_string"},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsMandatory,
		rabbitmq.WithPublishOptionsPersistentDelivery,
		rabbitmq.WithPublishOptionsExchange("events"),
	)

	if err != nil {
		return err
	}

	fmt.Println("Message published")

	return nil
}
