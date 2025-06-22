# Use the official Golang image to create a build artifact.
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy all the source code first
COPY . .

# Download dependencies and create the go.sum file based on the source code
RUN go mod tidy

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Expose port 8080 and start the application
EXPOSE 8080

CMD ["./main"] 