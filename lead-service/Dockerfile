FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o lead-service ./cmd/main.go
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/lead-service .
EXPOSE 8080
CMD ["./lead-service"]
