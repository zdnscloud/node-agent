/*

Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package main

import (
	"flag"
	"net"

	"github.com/zdnscloud/cement/log"

	pb "github.com/zdnscloud/node-agent/proto"
	"github.com/zdnscloud/node-agent/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var addr string
	flag.StringVar(&addr, "listen", ":80", "server listen address")
	flag.Parse()

	log.InitLogger(log.Debug)

	svr := server.NewServer()
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterNodeAgentServer(grpcServer, &svr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen failed:%s", err.Error())
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("run grpc server failed:%s", err.Error())
	}
}
