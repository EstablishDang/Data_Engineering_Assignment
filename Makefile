export GOOGLE_API_PATH=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GRPC_GATEWAY_VERSION)/third_party/googleapis/
export PROTOC_GEN_SWAGGER_PATH=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GRPC_GATEWAY_VERSION)/
GRPC_GATEWAY_VERSION := $(shell grep grpc-gateway go.mod | head -n 1 | rev | cut -d ' ' -f 1 | rev)

OUTPUT_BIN = app_mssql_server

all:
	@rm -f bin/$(OUTPUT_BIN)
	@rm -f mssql_server/proto/v1/*.pb.*go
	@echo "Building software"
	@protoc --proto_path=. --proto_path=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GRPC_GATEWAY_VERSION)/third_party/googleapis/ \
	 --proto_path=$(GOPATH)/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@$(GRPC_GATEWAY_VERSION)/ \
	 --swagger_out=allow_delete_body=true,logtostderr=true,disable_default_errors=true:. \
	 --go_out=. --go_opt=paths=source_relative \
	 --go-grpc_out=require_unimplemented_servers=false,paths=source_relative:. \
	 --grpc-gateway_out=logtostderr=true:.  mssql_server/proto/v1/mssql_server.proto
	@go build -v -o bin/$(OUTPUT_BIN) cmd/main.go

