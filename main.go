package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

var (
    ctx    context.Context
    cancel context.CancelFunc
    dep    *Dep
)

func init() {
    var err error
    ctx, cancel = context.WithCancel(context.Background())
    // Initialise dependencies.
    dep, err = InitDependencies(ctx)
    if err != nil {
        panic(err)
    }
    err = StartLogger()
    if err != nil {
        panic(err)
    }
}

func main() {
    fmt.Println("Starting HTTP Server!")

    go logUniqueRequests()
    go readKafkaQueue()
    defer logFile.Close()
    defer func() {
        if err := dep.producer.Close(); err != nil {
            fmt.Printf("Failed to close producer: %v", err)
        }
        if err := dep.consumer.Close(); err != nil {
            fmt.Printf("Failed to close consumer: %v", err)
        }
    }()

    // Create MUX router.
    router := mux.NewRouter()

    // Add API handler.
    router.HandleFunc("/api/verve/accept", VerveHandler).Methods("GET")

    log.Fatal(http.ListenAndServe(":8080", router))
}

// readKafkaQueue starts a consumer, that listens to the kafka producer.
func readKafkaQueue() {
    for {
        msg, err := dep.consumer.ReadMessage(ctx)
        if err != nil {
            fmt.Printf("Failed to read message: %v", err)
        }
        fmt.Printf("Received message at offset %d\n", msg.Offset)
        fmt.Printf("Unique requests in the minute: %s = %s\n", string(msg.Key), string(msg.Value))
    }
}
