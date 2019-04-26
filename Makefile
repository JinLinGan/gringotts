build:
	protoc -I proto proto/message.proto  --go_out=plugins=grpc:message

check-global-var:
	golangci-lint run --no-config  -E gochecknoglobals  ./...

check:
	golangci-lint run  ./...