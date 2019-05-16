package main

import (
	"context"
	"flag"
	"time"

	"github.com/zdnscloud/cement/log"
	"github.com/zdnscloud/node-agent/client"
	pb "github.com/zdnscloud/node-agent/proto"
)

func main() {
	var addr, path string
	flag.StringVar(&addr, "server", ":80", "server listen address")
	flag.StringVar(&path, "path", "/home/vagrant/workspace/code/go/src/github.com/zdnscloud/node-agent/vendor/google.golang.org/grpc", "path to get size")
	flag.Parse()

	log.InitLogger(log.Debug)
	defer log.CloseLogger()

	cli, err := client.NewClient(addr, 10*time.Second)
	if err != nil {
		log.Fatalf("connect to server failedd:%s", err.Error())
	}

	req := pb.GetDirectorySizeRequest{
		Path: path,
	}

	reply, err := cli.GetDirectorySize(context.TODO(), &req)
	if err != nil {
		log.Fatalf("get path size failed:%s", err.Error())
	} else {
		log.Infof("get size:%v", reply.Size)
	}
}
