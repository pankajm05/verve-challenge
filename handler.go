package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// PostData is the body for the post call made the endpoint that is provided as the query parameter.
type PostData struct {
    Count     int    `json:"count"`
    Timestamp string `json:"timestamp"`
}

// VerveHandler is handler for the api /api/verve/accept
func VerveHandler(w http.ResponseWriter, r *http.Request) {
    // parse request params
    idParam := getStringParam(r, "id")
    endpointParam := getStringParam(r, "endpoint")

    // Validate request params.
    id, errMsg := validateIdParam(idParam)
    if errMsg != nil {
        ErrorHandler(w, *errMsg, true)
        return
    }

    // store the ID for the given minute.
    _, err := dep.cacheStore.SetIDInCache(ctx, id)
    if err != nil {
        ErrorHandler(w, err.Error(), false)
        return
    }

    if endpointParam == "" {
        ResponseHandler(w, true)
        return
    }

    // Make the external API call in case endpoint is provided.
    var val int
    val, err = dep.cacheStore.GetUniqueIDsForCurrentMinute(ctx)
    if err != nil {
        ErrorHandler(w, err.Error(), false)
        return
    }
    err = externalRequestHandler(endpointParam, val)
    if err != nil {
        ErrorHandler(w, err.Error(), false)
        return
    }
    ResponseHandler(w, true)
    return
}

// externalRequestHandler makes an external API request.
func externalRequestHandler(endpoint string, uniqueCount int) error {
    postBody := PostData{
        Count:     uniqueCount,
        Timestamp: time.Now().UTC().Format(time.RFC3339),
    }
    jsonData, err := json.Marshal(postBody)
    if err != nil {
        fmt.Printf("Failed to marshal data: %v\n", err)
        return err
    }

    req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Printf("Failed to create request: %v\n", err)
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := dep.httpClient.Do(req)
    if err != nil {
        fmt.Printf("Failed to send request to endpoint: %v\n", err)
        return err
    }
    log.Printf("POST Request to %s with data: %s, returned status code: %d", endpoint, jsonData, resp.StatusCode)
    defer resp.Body.Close()
    return nil
}
