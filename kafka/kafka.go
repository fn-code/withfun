package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

// Config kafka config
type Config struct {
	Brokers    []string `mapstructure:"brokers" validate:"required"`
	GroupID    string   `mapstructure:"groupID" validate:"required,gte=0"`
	InitTopics bool     `mapstructure:"initTopics"`
}

// NewKafkaConn create new kafka connection
func NewKafkaConn(ctx context.Context, kafkaCfg *Config) (*kafka.Conn, error) {
	return kafka.DialContext(ctx, "tcp", kafkaCfg.Brokers[0])
}
