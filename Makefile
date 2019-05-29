build:
	protoc -I proto proto/message.proto  --go_out=plugins=grpc:message

check-all:
	golangci-lint run --no-config  --enable-all  ./...

check:
	golangci-lint run  ./...

tidy:
	GO111MODULE=on  go mod tidy
	GO111MODULE=on  go mod vendor

run-agent:
	go run -race ./gringotts-agent/main.go start

test:
	go test -race ./...