test:
	go test -v -cover ./...

server:
	go run cmd/app/main.go

mock:
	mockgen -package mocks -destination internal/usecase/mocks/httpclient.go github.com/zura-t/go_delivery_system/internal/usecase/httpclient HttpClientI

.PHONY: test server