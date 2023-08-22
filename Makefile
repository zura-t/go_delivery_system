test:
	go test -v -cover ./...

server:
	go run main.go

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
		--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=delivery_system \
    proto/*.proto
		statik -src=./doc/swagger -dest=./doc

.PHONY: test server proto