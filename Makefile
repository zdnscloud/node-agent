REGISTRY_NAME = zdnscloud/node-agent
IMAGE_VERSION = latest

BRANCH=`git branch | sed -n '/\* /s///p'`
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.build=${BUILD}"

all: grpc

grpc: proto/nodeagent.pb.go

proto/nodeagent.pb.go: proto/nodeagent.proto
	cd proto && protoc -I/usr/local/include -I. --go_out=plugins=grpc:. nodeagent.proto

clean:
	rm -f proto/nodeagent.pb.go

container:
	go mod vendor
	#docker build -t $(REGISTRY_NAME):${BRANCH} --build-arg version=${VERSION} --build-arg buildtime=${BUILD} --no-cache .
	docker build -t $(REGISTRY_NAME):${IMAGE_VERSION} --build-arg version=${VERSION} --build-arg buildtime=${BUILD} --no-cache .
	docker image prune -f
	rm -fr vendor

build:
	CGO_ENABLED=0 GOOS=linux go build ${LDFLAGS}

.PHONY: all clean
