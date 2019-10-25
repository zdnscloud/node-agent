REGISTRY_NAME = zdnscloud/node-agent
IMAGE_VERSION = v1.3

all: grpc

grpc: proto/nodeagent.pb.go

proto/nodeagent.pb.go: proto/nodeagent.proto
	cd proto && protoc -I/usr/local/include -I. --go_out=plugins=grpc:. nodeagent.proto

clean:
	rm -f proto/nodeagent.pb.go

container:
	docker build -t $(REGISTRY_NAME):$(IMAGE_VERSION) ./ --no-cache

.PHONY: all clean
