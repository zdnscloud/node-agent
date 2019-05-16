all: grpc

grpc: proto/nodeagent.pb.go

proto/nodeagent.pb.go: proto/nodeagent.proto
	cd proto && protoc -I/usr/local/include -I. --go_out=plugins=grpc:. nodeagent.proto

clean:
	rm -f proto/nodeagent.pb.go

.PHONY: all clean
