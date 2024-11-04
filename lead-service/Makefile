.PHONY: proto clean run test build generate

proto:
	protoc --proto_path=proto \
		--go_out=proto/pb --go_opt=paths=source_relative \
		--go-grpc_out=proto/pb --go-grpc_opt=paths=source_relative \
		proto/*.proto

clean:
	rm -f proto/pb/*.pb.go

run:
	go run cmd/main.go

test:
	go test ./... -v

build:
	go build -o bin/lead-service cmd/main.go

generate: clean proto