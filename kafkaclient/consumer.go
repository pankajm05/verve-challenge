package kafkaclient

import (
    "context"

    "github.com/segmentio/kafka-go"
)

// Consumer is a Kafka consumer structure.
type Consumer struct {
    reader *kafka.Reader
}

// NewConsumer creates a new consumer instance.
func NewConsumer(cfg Config, cCfg ConsumerConfig) *Consumer {
    return &Consumer{
        reader: kafka.NewReader(kafka.ReaderConfig{
            Brokers:   []string{cfg.Addr},
            Topic:     cfg.Topic,
            GroupID:   cCfg.GroupID,
            Partition: cCfg.Partition,
            MinBytes:  cCfg.MinBytes,
            MaxBytes:  cCfg.MaxBytes,
        }),
    }
}

// ReadMessage reads a message from the Kafka topic.
func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
    return c.reader.ReadMessage(ctx)
}

// Close closes the consumer reader.
func (c *Consumer) Close() error {
    return c.reader.Close()
}
