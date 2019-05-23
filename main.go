package main

import (
	"flag"
	"net"

	"github.com/zdnscloud/cement/log"

	"github.com/zdnscloud/node-agent/clusteragent"
	pb "github.com/zdnscloud/node-agent/proto"
	"github.com/zdnscloud/node-agent/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var addr string
	var clusterAgentAddr string
	var nodeName string

	flag.StringVar(&addr, "listen", ":80", "server listen address")
	flag.StringVar(&clusterAgentAddr, "cluster", "cluster-agent.zcloud.svc", "server listen address")
	flag.StringVar(&nodeName, "node", "", "node name this node agent resides")
	flag.Parse()

	log.InitLogger(log.Debug)

	if nodeName == "" {
		log.Fatalf("node name isn't specified")
	} else if clusterAgentAddr == "" {
		log.Fatalf("cluster agent addr isn't specified")
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen failed:%s", err.Error())
	}

	clusteragent.StartHeartbeat(clusterAgentAddr, nodeName, addr)

	svr := server.NewServer()
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterNodeAgentServer(grpcServer, &svr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("run grpc server failed:%s", err.Error())
	}
}
