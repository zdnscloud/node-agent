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
	addr := "127.0.0.1:9001"
	iqn := "iqn.1992-08.com.netapp:sn.142255150"
	cli, err := client.NewClient(addr, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}
	blocks, err := cli.IscsiGetBlocks(ctx, &pb.IscsiGetBlocksRequest{
		Iqn: iqn,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	var devices []string
	for _, info := range blocks.IscsiBlock {
		fmt.Println(info.Blocks)
		if len(info.Blocks) > 0 {
			path, err := cli.IscsiGetMultipaths(ctx, &pb.IscsiGetMultipathsRequest{
				Devs: info.Blocks,
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			devices = append(devices, multipathDir+path.Dev)
		} else {
			for _, dev := range info.Blocks {
				devices = append(devices, deviceDir+dev)
			}
		}
	}
	fmt.Println(devices)
}
