package kafkaclient

import (
    "context"
    "time"

    "github.com/segmentio/kafka-go"
)

const (
    maxRetries = 5
    increment  = 10 * time.Second
)

// Producer is a Kafka producer structure.
type Producer struct {
    Writer *kafka.Writer
}

// NewProducer creates a new producer instance.
func NewProducer(cfg Config, pCfg ProducerConfig) *Producer {
    return &Producer{
        Writer: &kafka.Writer{
            Addr:                   kafka.TCP(cfg.Addr),
            Topic:                  cfg.Topic,
            Balancer:               &kafka.LeastBytes{},
            RequiredAcks:           pCfg.RequiredAcks,
            Async:                  pCfg.Async,
            AllowAutoTopicCreation: pCfg.AutoTopicCreation,
        },
    }
}

// SendMessage sends a message to the Kafka topic.
func (p *Producer) SendMessage(ctx context.Context, key, value []byte) error {
    var (
        err        error
        retryDelay time.Duration
    )
    for i := 0; i < maxRetries; i++ {
        retryDelay += increment
        err = p.Writer.WriteMessages(ctx, kafka.Message{
            Key:   key,
            Value: value,
        })
        if err == nil {
            return nil
        }
        time.Sleep(retryDelay)

    }
    return err
}

// Close closes the producer Writer.
func (p *Producer) Close() error {
    return p.Writer.Close()
}
