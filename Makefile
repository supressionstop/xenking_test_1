.PHONY: lint
lint:
	golangci-lint run --config ./config/golangci.yml

.PHONY: tests
tests:
	go test ./internal/...

.PHONY: run
run:
	docker compose up -d --build

.PHONY: stop
stop:
	docker compose stop

# helpers

.PHONY: proto
proto:
	protoc --proto_path=api \
			--go_out=. \
			--go-grpc_out=. \
			lines.proto

.PHONY: mock
mock:
	mockery --config ./config/mockery.yml

.PHONY: debug-run
debug-run:
	docker compose -f docker-compose.yml -f docker-compose.debug.yml up --build

grpc-client:
	go build -o tools/grpc_client github.com/supressionstop/xenking_test_1/cmd/grpc_client