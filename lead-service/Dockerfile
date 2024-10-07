# Build stage
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the Go binary and place it in the /app directory
RUN go build -o lead-service ./cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
# Copy the binary from the builder stage to the current directory
COPY --from=builder /app/lead-service /app/lead-service
# Ensure the binary has execution permissions
RUN chmod +x /app/lead-service
EXPOSE 8080
# Use the full path to execute the binary
CMD ["/app/lead-service"]
