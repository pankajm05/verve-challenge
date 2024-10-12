package main

import (
    "log"
    "net/http"
    "strconv"
    "strings"
)

func validateIdParam(param string) (int64, *string) {
    var errMsg string
    if param == "" {
        errMsg = "'id' is required query parameter"
        return -1, &errMsg
    }

    id, err := strconv.ParseInt(strings.TrimSpace(param), 10, 64)
    if err != nil {
        errMsg = "'id' query parameter must be an integer."
        return -1, &errMsg
    }

    // Check for non-negative integers.
    if id <= 0 {
        errMsg = "'id' query parameter must be greater than or equal to 0."
        return -1, &errMsg
    }

    return id, nil
}

func getStringParam(r *http.Request, key string) string {
    param := r.URL.Query().Get(key)
    return strings.Trim(param, " \"")
}

// ResponseHandler is the common method to write API response.
func ResponseHandler(w http.ResponseWriter, requestSuccessful bool) {
    result := "ok"
    if !requestSuccessful {
        result = "failed"
    }
    _, err := w.Write([]byte(result))
    if err != nil {
        log.Fatal(w)
    }
}

// ErrorHandler is the common method to handle error.
func ErrorHandler(w http.ResponseWriter, errMsg string, isClientError bool) {
    errorCode := http.StatusInternalServerError
    if isClientError {
        errorCode = http.StatusBadRequest
    }
    http.Error(w, errMsg, errorCode)
    ResponseHandler(w, false)
}
