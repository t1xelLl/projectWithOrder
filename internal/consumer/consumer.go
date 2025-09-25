package consumer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/t1xelLl/projectWithOrder/internal/entities"
	"github.com/t1xelLl/projectWithOrder/internal/service"
)

func StartConsumer(service *service.Service, reader *kafka.Reader) {
	go func() {
		for {
			message, err := reader.ReadMessage(context.Background())
			if err != nil {
				logrus.Printf("kafka read error :%v\n", err)
				continue
			}

			var order entities.Order
			if err := json.Unmarshal(message.Value, &order); err != nil {
				logrus.Printf("Invalid message :%v\n", err)
			}

			if err := service.CreateOrder(context.Background(), &order); err != nil {
				logrus.Printf("failed to save order :%v\n", err)
			} else {
				logrus.Printf("order saved :%v\n", order.OrderUID)
			}

		}
	}()
}
