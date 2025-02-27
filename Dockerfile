# syntax=docker/dockerfile:1

# --------------------------
# 1) Build Stage
# --------------------------
FROM golang:1.24.0-alpine AS builder

# Create and set our working directory
WORKDIR /app

# Copy go.mod and go.sum first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the app source
COPY . .

# Build the Go app and name the output binary 'main'
RUN go build -o main .

# --------------------------
# 2) Final Stage
# --------------------------
FROM alpine:3.18

# Create and set our working directory
WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=builder /app/main /app/main

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./main"]