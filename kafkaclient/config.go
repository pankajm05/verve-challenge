package kafkaclient

import (
    "github.com/segmentio/kafka-go"
)

// Config common configuration used by both producers and consumers.
type Config struct {
    Addr  string
    Topic string
}

// ProducerConfig specific configuration for producers.
type ProducerConfig struct {
    Async             bool
    AutoTopicCreation bool
    RequiredAcks      kafka.RequiredAcks
}

// ConsumerConfig specific configuration for consumers.
type ConsumerConfig struct {
    GroupID        string
    Partition      int
    MinBytes       int
    MaxBytes       int
    CommitInterval int // In milliseconds
    StartOffset    int64
}
