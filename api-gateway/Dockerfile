# Build stage
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the Go binary and place it directly in the /app directory
RUN go build -o /app/api-gateway ./cmd/main.go

# Final stage
FROM alpine:latest
WORKDIR /root/
# Copy the binary from the builder stage to the current directory
COPY --from=builder /app/api-gateway .
# Ensure the binary has execution permissions
RUN chmod +x api-gateway
EXPOSE 8080
CMD ["./api-gateway"]
