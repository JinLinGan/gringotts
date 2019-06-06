DOCKERPATH ?= /private/var/gringotts/docker

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

run-server:
	go run -race ./gringotts-server/main.go start

run-mysql:
	docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -v ${DOCKERPATH}/mysql/data:/var/lib/mysql  mysql:5

stop-mysql:
	docker stop mysql

rm-mysql:
	docker rm mysql