package main

import (
	"context"
	"fmt"
	nodeagaentclient "github.com/zdnscloud/node-agent/client"
	pb "github.com/zdnscloud/node-agent/proto"
	"time"
)

func main() {
	ctx := context.TODO()
	addr := "10.0.0.146:8888"
	cli, err := nodeagaentclient.NewClient(addr, 10*time.Second)
	if err != nil {
		fmt.Println("1111", err)
		return
	}
	if _, err := cli.IscsiDiscovery(ctx, &pb.IscsiDiscoveryRequest{
		Host: "10.0.0.104",
		Port: "3260",
		Iqn:  "iqn.1992-08.com.netapp:sn.142255150",
	}); err != nil {
		fmt.Println("2222", err)
		return
	}
	if _, err := cli.IscsiLogin(ctx, &pb.IscsiLoginRequest{
		Host: "10.0.0.104",
		Port: "3260",
		Iqn:  "iqn.1992-08.com.netapp:sn.142255150",
	}); err != nil {
		fmt.Println("3333", err)
		return
	}
}
