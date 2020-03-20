DOCKERPATH ?= /private/var/gringotts/docker

build-grpc: proto/message.proto
	protoc -I proto proto/message.proto  --go_out=plugins=grpc:pkg/message

check-all:
	golangci-lint run --no-config  --enable-all  ./...

check:
	golangci-lint run  ./...

tidy:
	GO111MODULE=on  go mod tidy
	GO111MODULE=on  go mod vendor

run-agent: build-grpc
	go run -race ./gringotts-agent/main.go start

run-server: build-grpc
	go run -race ./gringotts-server/main.go start

build-agent-linux: build-grpc
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/agent ./gringotts-agent/main.go

build-server-linux: build-grpc
	GOOS=linux GOARCH=amd64 go build -gcflags "all=-N -l" -o build/bin/server ./gringotts-server/main.go

build-dev-image:
	docker build  -f ./scripts/dockerfiles/dlv/Dockerfile . -t harbor.gk8s.ete.ffcs.cn/gringotts/devtool:latest
create-docker-network:
	-docker network create gringotts

run-agent-docker: build-agent-linux create-docker-network
	-docker rm -f agent
	docker run \
		--entrypoint /gringotts/bin/agent \
		-v /Users/jinlin/code/golang/src/github.com/jinlingan/gringotts/build/bin:/gringotts/bin \
		--name agent \
		--network gringotts \
		harbor.gk8s.ete.ffcs.cn/gringotts/devtool:latest start

run-debug-agent-docker: build-agent-linux create-docker-network
	-docker rm -f agent
	docker run -d --entrypoint dlv \
  		--network gringotts \
		-p 2345:2345 \
 		-v /Users/jinlin/code/golang/src/github.com/jinlingan/gringotts/build/bin:/gringotts/bin  \
 		--name agent harbor.gk8s.ete.ffcs.cn/gringotts/devtool:latest  \
 		--listen=:2345 --headless=true --api-version=2 --accept-multiclient exec /gringotts/bin/agent start

run-debug-server-docker: build-server-linux create-docker-network
	-docker rm -f server
	docker run -d --entrypoint dlv \
		--network gringotts \
		-p 2346:2346 \
  		-v /Users/jinlin/code/golang/src/github.com/jinlingan/gringotts/build/bin:/gringotts/bin  \
  		--name server harbor.gk8s.ete.ffcs.cn/gringotts/devtool:latest  \
  		--listen=:2346 --headless=true --api-version=2 --accept-multiclient exec /gringotts/bin/server start

init-run-local-mysql: create-docker-network
	-docker rm -f mysql
	docker run -d \
    	--network gringotts \
    	-p 3306:3306 \
    	-e MYSQL_ROOT_PASSWORD=root.123 \
      	-v /Users/jinlin/code/golang/src/github.com/jinlingan/gringotts/scripts/database/init:/docker-entrypoint-initdb.d  \
      	--name mysql mysql:5