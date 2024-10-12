# Verve Tech Challenge

This is my submission to the Verve Tech challenge.

The service supports a single GET endpoint:

#### Base URL:

`/api/verve/accept`

#### Query Parameters:

|   Name   |  type   | Required | Example                                                |
|:--------:|:-------:|:--------:|:-------------------------------------------------------|
|    id    | integer |   true   | `/api/verve/accept?id=2`                               |
| endpoint | string  | optional | `/api/verve/accept?id=10&endpoint=https://example.com` |

#### How to run:

Docker Image Link: https://hub.docker.com/r/pankajm05/verve-go-app


```
// clone the repository.
// Assuming you are in the root directory of the repo.
> docker-compose build
> docker-compose up
// Wait for the message: "Starting HTTP Server!" in the verve-go-app service
// The API should be accessible on the url: http://localhost:8080/api/verve/accept
```
