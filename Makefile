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

.PHONY: generate
generate: proto mock

.PHONY: debug-run
debug-run:
	docker compose -f docker-compose.yml -f docker-compose.debug.yml up --build

.PHONY: grpc-client-run
grpc-client-run:
	go run github.com/supressionstop/xenking_test_1/cmd/grpc_client

.PHONY: lint-fix
lint-fix:
	golangci-lint run --config ./config/golangci.yml --fix

.PHONY: run-local
run-local:
	APP_ENV=local APP_NAME=applocal go run github.com/supressionstop/xenking_test_1/cmd/processor