package main

import (
    "fmt"
    "log"
    "os"
    "strconv"
    "time"
)

var logFile *os.File

// logUniqueRequests reports the number of unique request in the elapsed minute.
func logUniqueRequests() {
    tick := time.Tick(time.Minute)
    for {
        select {
        case <-tick:
            prevMintue := time.Now().UTC().Add(-1 * time.Minute).Truncate(time.Minute)
            val, err := dep.cacheStore.GetUniqueIDsForPreviousMinute(ctx)
            if err != nil {
                log.Printf("Error getting unique ids from cache: %v", err)
                continue
            }
            err = dep.producer.SendMessage(ctx, []byte(prevMintue.Format("15:04")), []byte(strconv.Itoa(val)))
            if err != nil {
                fmt.Printf("Failed to send message: %v", err)
            }
        }
    }
}

// StartLogger starts logging to a file.
func StartLogger() error {
    var err error
    fmt.Println("Initializing logger")
    logFile, err = os.OpenFile("request.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

    if err != nil {
        log.Fatalln("Failed to open log file:", err)
        return err
    }

    log.SetFlags(log.LstdFlags | log.LUTC)
    log.SetOutput(logFile)
    return nil
}
