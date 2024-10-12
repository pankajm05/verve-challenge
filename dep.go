package main

import (
    "context"
    "fmt"
    "net/http"

    "verve-tech-challenge/kafkaclient"
    "verve-tech-challenge/redisclient"

    "github.com/segmentio/kafka-go"
)

type Dep struct {
    httpClient *http.Client
    cacheStore *redisclient.Cache
    producer   *kafkaclient.Producer
    consumer   *kafkaclient.Consumer
}

func InitDependencies(ctx context.Context) (*Dep, error) {
    // Initialise Redis Cache.
    var d Dep
    var err error
    d.httpClient = &http.Client{}
    d.cacheStore, err = redisclient.NewCache(ctx, redisclient.Config{
        Addr:     "redis:6379",
        Password: "",
        Database: 0,
    }, true)
    if err != nil {
        return nil, err
    }
    fmt.Println("Initialising Kafka Queue")
    // Define common Kafka configuration
    cfg := kafkaclient.Config{
        Addr:  "kafka1:9092",
        Topic: kafkaTopic,
    }

    // Define producer-specific configuration
    pCfg := kafkaclient.ProducerConfig{
        Async:             false,
        AutoTopicCreation: true,
        RequiredAcks:      kafka.RequireAll,
    }

    // Initialize the producer
    d.producer = kafkaclient.NewProducer(cfg, pCfg)
    fmt.Println("Initialising Kafka producer")
    err = d.producer.SendMessage(ctx, []byte("key"), []byte("Hello, Kafka!"))
    if err != nil {
        return nil, err
    }

    // Define consumer-specific configuration
    cCfg := kafkaclient.ConsumerConfig{
        GroupID:        kafkaConsGroup,
        Partition:      0,
        MinBytes:       10e3, // 10KB
        MaxBytes:       10e6, // 10MB
        CommitInterval: 1000, // 1 second
        StartOffset:    kafka.FirstOffset,
    }

    // Initialize the consumer
    fmt.Println("Initialising Kafka producer")
    d.consumer = kafkaclient.NewConsumer(cfg, cCfg)
    return &d, nil
}
