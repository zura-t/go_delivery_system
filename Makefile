test:
	go test -v -cover ./...

server:
	go run cmd/app/main.go

.PHONY: test server