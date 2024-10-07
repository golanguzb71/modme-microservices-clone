FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api-gateway ./cmd/main.go
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]
