.PHONY: lint
lint:
	echo "todo linters"

.PHONY: tests
tests:
	echo "todo tests"

.PHONY: run
run:
	echo "todo run"

.PHONY: stop
stop:
	echo "todo stop"

### helpers

.PHONY: proto
proto:
	protoc --proto_path=api \
			--go_out=. \
			--go-grpc_out=. \
			lines.proto