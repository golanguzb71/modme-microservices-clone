FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o lead-service ./cmd/main.go

EXPOSE 8080

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/lead-service ./
COPY --from=builder /app/config/config.yaml ./config/

CMD ["./lead-service"]
