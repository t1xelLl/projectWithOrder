package kafka

import "github.com/segmentio/kafka-go"

type Consumer struct {
	reader *kafka.Reader
}
