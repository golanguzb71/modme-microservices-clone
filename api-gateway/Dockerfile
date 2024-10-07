# Build stage
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the Go binary with static linking to ensure compatibility with Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
# Copy the binary from the builder stage to the current directory
COPY --from=builder /app/api-gateway /app/api-gateway
# Ensure the binary has execution permissions
RUN chmod +x /app/api-gateway
# Verify that the binary exists
RUN ls -l /app/api-gateway
EXPOSE 8080
CMD ["/app/api-gateway"]
