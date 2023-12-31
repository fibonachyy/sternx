PROTO_INCLUDE ?= /usr/include # set the location of protocol buffer libs
PROTO_OUT ?= internal/api

# Generate certificates
.PHONY: generate-certs
generate-certs:
	openssl genpkey -algorithm RSA -out certs/keyFile.pem
	openssl req -new -key certs/keyFile.pem -out certs/csr.pem
	openssl x509 -req -in certs/csr.pem -signkey certs/keyFile.pem -out certs/certFile.pem
	rm certs/csr.pem

# Generate gRPC code
.PHONY: generate-proto
generate-proto:
	rm -f internal/api/*.go
	rm -f doc/swagger/*.swagger.json
	protoc -I $(PROTO_INCLUDE) --proto_path=proto --go_out=$(PROTO_OUT) --go_opt=paths=source_relative \
        --go-grpc_out=$(PROTO_OUT) --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=$(PROTO_OUT) --grpc-gateway_opt=paths=source_relative \
        proto/*.proto

	# Generate OpenAPI JSON separately
	protoc -I $(PROTO_INCLUDE) --proto_path=proto \
        --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=user_service \
        proto/*.proto

.PHONY: start
start:
	go run main.go server --config ./config.yaml
	
.PHONY: all
all: generate-certs generate-proto

