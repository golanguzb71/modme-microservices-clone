# Build stage
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the Go binary and place it in the /app directory
RUN go build -o api-gateway ./cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
# Copy the binary from the builder stage to the current directory
COPY --from=builder /app/api-gateway /app/api-gateway
# Ensure the binary has execution permissions
RUN chmod +x /app/api-gateway
EXPOSE 8080
# Use the full path to execute the binary
CMD ["/app/api-gateway"]
