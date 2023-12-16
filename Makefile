PROTO_INCLUDE ?= /usr/include

# Install recursive submodules
.PHONY: install-submodules
install-submodules:
	git submodule update --init --recursive

# Generate certificates
.PHONY: generate-certs
generate-certs:
	openssl genpkey -algorithm RSA -out certs/keyFile.pem
	openssl req -new -key certs/keyFile.pem -out certs/csr.pem
	openssl x509 -req -in certs/csr.pem -signkey certs/keyFile.pem -out certs/certFile.pem
	rm certs/csr.pem

# Generate gRPC code
.PHONY: generate-user
generate-user:
	protoc -I . -I ./@googleapis -I $(PROTO_INCLUDE) -I ./internal/api/user/ \
		--go-grpc_out=paths=source_relative:. \
  		--go_out=paths=source_relative:. \
		./internal/api/user/user.proto

# Generate gRPC-gateway code
.PHONY: generate-user-gateway
generate-user-gateway:
	protoc -I ./@googleapis -I $(PROTO_INCLUDE) -I ./internal/api/user/ -I . \
		--grpc-gateway_out=logtostderr=true,paths=source_relative:./internal/api/user/ \
		--openapiv2_out=./internal/api/user/ \
		./internal/api/user/user.proto

# Default target, installs submodules and generates certificates and code
.PHONY: all
all: install-submodules generate-certs generate-user generate-user-gateway