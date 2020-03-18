package main

import (
	"context"
	"fmt"
	"github.com/zdnscloud/node-agent/client"
	pb "github.com/zdnscloud/node-agent/proto"
	"time"
)

const (
	multipathDir = "/dev/mapper/"
	deviceDir    = "/dev/"
)

var ctx = context.TODO()

func main() {
	addr := "10.0.0.147:9001"
	cli, err := client.NewClient(addr, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	req := pb.GetDisksInfoRequest{}
	reply, err := cli.GetDisksInfo(context.TODO(), &req)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reply.Infos)
}
