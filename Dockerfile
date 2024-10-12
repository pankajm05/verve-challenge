# Use Golang image to build the application
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download dependencies and build the application for Linux
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Use a smaller image for running the built app
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .

# Expose the port that the service will run on
EXPOSE 8080

# Command to run the application
CMD ["./app"]
