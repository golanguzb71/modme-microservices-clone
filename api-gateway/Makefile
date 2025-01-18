.PHONY: proto clean run test swagger

proto:
	protoc --proto_path=grpc/proto \
		--go_out=grpc/proto/pb --go_opt=paths=source_relative \
		--go-grpc_out=grpc/proto/pb --go-grpc_opt=paths=source_relative \
		grpc/proto/*.proto

clean:
	rm -f grpc/proto/pb/*.pb.go
	rm -f docs/docs.go docs/swagger.json docs/swagger.yaml

clean_omit:
	find ./grpc/proto/pb -name "*.pb.go" -exec sh -c 'sed -i "s/,omitempty//" {}' \;
run:
	go run cmd/main.go

test:
	go test ./... -v

swagger:
	swag init -g cmd/main.go

build:
	go build -o bin/lead-service cmd/main.go

generate: clean proto swagger clean_omit
