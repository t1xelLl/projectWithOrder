package kafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/t1xelLl/projectWithOrder/configs"
)

func NewKafkaReader(cfg configs.Kafka) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{cfg.Host + ":" + cfg.Port},
		Topic:       cfg.Topic,
		GroupID:     cfg.GroupID,
		StartOffset: kafka.LastOffset,
	})
}
