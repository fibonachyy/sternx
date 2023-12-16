.PHONY: generate-certs

generate-certs:
	openssl genpkey -algorithm RSA -out certs/keyFile.pem
	openssl req -new -key certs/keyFile.pem -out certs/csr.pem
	openssl x509 -req -in certs/csr.pem -signkey certs/keyFile.pem -out certs/certFile.pem
	rm certs/csr.pem

.PHONY: generate-proto
generate-proto:
	protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. internal/api/**/*.proto

.PHONY: generate-user-gateway
generate-user-gateway:
	protoc -I ./protobuf -I ./internal/api/user/ -I . --grpc-gateway_out=logtostderr=true,paths=source_relative:./internal/api/user/ --openapiv2_out=./internal/api/user/ ./internal/api/user/user.proto
	 