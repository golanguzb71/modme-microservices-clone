.PHONY: proto clean clean_omit run test swagger build generate check_deps

PROTO_SRC = grpc/proto
PROTO_OUT = grpc/proto/pb
SWAGGER_OUT = docs
CMD_ENTRY = cmd/main.go
BIN_NAME = bin/lead-service

check_deps:
	@command -v protoc > /dev/null || (echo "protoc is not installed"; exit 1)
	@command -v swag > /dev/null || (echo "swag is not installed"; exit 1)

proto:
	protoc --proto_path=$(PROTO_SRC) \
		--go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_SRC)/*.proto

clean:
	rm -f $(PROTO_OUT)/*.pb.go
	rm -f $(SWAGGER_OUT)/docs.go $(SWAGGER_OUT)/swagger.json $(SWAGGER_OUT)/swagger.yaml

clean_omit:
	find $(PROTO_OUT) -name "*.pb.go" -exec sed -i "s/,omitempty//g" {} +

run:
	go run $(CMD_ENTRY)

test:
	go test ./... -v

swagger:
	swag init -g $(CMD_ENTRY)

build:
	go build -o $(BIN_NAME) $(CMD_ENTRY)

generate: check_deps clean proto swagger clean_omit
