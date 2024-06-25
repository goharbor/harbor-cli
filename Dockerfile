# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

RUN go build -o harbor ./cmd/harbor

# Stage 2: Create a small image for the Go binary
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/harbor .

CMD ["./harbor"]