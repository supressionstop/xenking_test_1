FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app
COPY go.mod go.mod
RUN go mod download
COPY . .
RUN go build -v -o ./processor ./cmd/processor

FROM alpine
WORKDIR /app
COPY config config
COPY --from=builder /app/processor ./processor
ENTRYPOINT ./processor